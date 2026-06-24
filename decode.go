package fastabi

import (
	"fmt"
	"math/big"
)

func Decode(ts []ParamType, data []byte) ([]any, error) {
	return decodeTuple(ts, data, 0)
}

func Decode1(t ParamType, data []byte) (any, error) {
	vals, err := decodeTuple([]ParamType{t}, data, 0)
	if err != nil {
		return nil, err
	}
	return vals[0], nil
}

func decodeTuple(ts []ParamType, data []byte, offset int) ([]any, error) {
	vals := make([]any, len(ts))
	pos := offset

	for i, t := range ts {
		if t.IsDynamic() {
			if pos+32 > len(data) {
				return nil, fmt.Errorf("decode tuple[%d]: offset overrun", i)
			}
			off := ReadUint64At(data, pos)
			pos += 32
			v, err := decodeDynamic(t, data, offset+int(off))
			if err != nil {
				return nil, fmt.Errorf("decode tuple[%d] %s: %w", i, t.Kind, err)
			}
			vals[i] = v
		} else {
			sz := t.StaticSize()
			if pos+sz > len(data) {
				return nil, fmt.Errorf("decode tuple[%d]: static overrun", i)
			}
			v, err := decodeStatic(t, data, pos)
			if err != nil {
				return nil, fmt.Errorf("decode tuple[%d] %s: %w", i, t.Kind, err)
			}
			pos += sz
			vals[i] = v
		}
	}
	return vals, nil
}

func decodeStatic(t ParamType, data []byte, offset int) (any, error) {
	if offset+32 > len(data) {
		return nil, fmt.Errorf("static decode overrun at %d", offset)
	}
	slot := data[offset : offset+32]

	switch t.Kind {
	case KindUint256:
		var u U256
		u.SetBytes(slot)
		return &u, nil
	case KindUint128:
		return new(big.Int).SetBytes(slot[16:]), nil
	case KindUint64:
		return readUint64BE(slot[24:]), nil
	case KindUint32:
		return uint64(readUint32BE(slot[28:])), nil
	case KindUint24:
		return uint64(slot[29])<<16 | uint64(slot[30])<<8 | uint64(slot[31]), nil
	case KindUint16:
		return uint64(slot[30])<<8 | uint64(slot[31]), nil
	case KindUint8:
		return uint64(slot[31]), nil
	case KindInt256:
		return decodeInt256(slot), nil
	case KindInt128:
		return decodeInt128(slot), nil
	case KindInt64:
		return int64(readUint64BE(slot[24:])), nil
	case KindInt32:
		return int64(int32(readUint32BE(slot[28:]))), nil
	case KindInt24:
		v := int32(slot[29])<<16 | int32(slot[30])<<8 | int32(slot[31])
		// Sign extend from 24-bit
		if v&0x800000 != 0 {
			v |= ^0xFFFFFF
		}
		return int64(v), nil
	case KindInt16:
		return int64(int16(slot[30])<<8 | int16(slot[31])), nil
	case KindInt8:
		return int64(int8(slot[31])), nil
	case KindAddress:
		var addr Address
		copy(addr[:], slot[12:])
		return addr, nil
	case KindBool:
		return slot[31] != 0, nil
	case KindBytes32:
		var b [32]byte
		copy(b[:], slot)
		return b, nil
	case KindFixedArr:
		return decodeFixedArray(t, data, offset)
	case KindTuple:
		return decodeTuple(t.TupleEl, data, offset)
	default:
		return nil, nil
	}
}

func decodeDynamic(t ParamType, data []byte, offset int) (any, error) {
	switch t.Kind {
	case KindBytes:
		return decodeDynBytes(data, offset), nil
	case KindString:
		b := decodeDynBytes(data, offset)
		return string(b), nil
	case KindArray:
		return decodeDynArray(t, data, offset)
	case KindFixedArr:
		return decodeFixedArray(t, data, offset)
	case KindTuple:
		return decodeTuple(t.TupleEl, data, offset)
	default:
		return nil, fmt.Errorf("unknown dynamic kind: %s", t.Kind)
	}
}

func decodeDynBytes(data []byte, offset int) []byte {
	if offset+32 > len(data) {
		return nil
	}
	length := int(ReadUint64At(data, offset))
	offset += 32
	if offset+length > len(data) {
		return nil
	}
	result := make([]byte, length)
	copy(result, data[offset:offset+length])
	return result
}

func decodeDynArray(t ParamType, data []byte, offset int) ([]any, error) {
	if offset+32 > len(data) {
		return nil, fmt.Errorf("array length overrun")
	}
	length := int(ReadUint64At(data, offset))
	offset += 32

	vals := make([]any, length)
	if t.Elem.IsDynamic() {
		// Read offsets first
		if offset+length*32 > len(data) {
			return nil, fmt.Errorf("array offsets overrun")
		}
		startOffset := offset - 32
		offsets := make([]int, length)
		for i := 0; i < length; i++ {
			offsets[i] = int(ReadUint64At(data, offset+i*32))
		}
		for i, off := range offsets {
			v, err := decodeDynamic(*t.Elem, data, startOffset+off)
			if err != nil {
				return nil, fmt.Errorf("array[%d]: %w", i, err)
			}
			vals[i] = v
		}
	} else {
		for i := 0; i < length; i++ {
			v, err := decodeStatic(*t.Elem, data, offset+i*32)
			if err != nil {
				return nil, fmt.Errorf("array[%d]: %w", i, err)
			}
			vals[i] = v
		}
	}
	return vals, nil
}

func decodeFixedArray(t ParamType, data []byte, offset int) ([]any, error) {
	n := t.Size
	vals := make([]any, n)
	startOffset := offset
	for i := 0; i < n; i++ {
		if t.Elem.IsDynamic() {
			if offset+32 > len(data) {
				return nil, fmt.Errorf("fixedArray[%d] offset overrun", i)
			}
			off := int(ReadUint64At(data, offset))
			offset += 32
			v, err := decodeDynamic(*t.Elem, data, startOffset+off)
			if err != nil {
				return nil, fmt.Errorf("fixedArray[%d]: %w", i, err)
			}
			vals[i] = v
		} else {
			v, err := decodeStatic(*t.Elem, data, offset)
			if err != nil {
				return nil, fmt.Errorf("fixedArray[%d]: %w", i, err)
			}
			offset += 32
			vals[i] = v
		}
	}
	return vals, nil
}

func decodeInt256(slot []byte) *big.Int {
	if slot[0]&0x80 == 0 {
		return new(big.Int).SetBytes(slot)
	}
	// Negative: two's complement
	inv := make([]byte, 32)
	for i := range slot {
		inv[i] = ^slot[i]
	}
	v := new(big.Int).SetBytes(inv)
	v.Add(v, big.NewInt(1))
	return v.Neg(v)
}

func decodeInt128(slot []byte) *big.Int {
	if slot[16]&0x80 == 0 {
		return new(big.Int).SetBytes(slot[16:])
	}
	inv := make([]byte, 16)
	for i := range slot[16:] {
		inv[i] = ^slot[16+i]
	}
	v := new(big.Int).SetBytes(inv)
	v.Add(v, big.NewInt(1))
	return v.Neg(v)
}

func ReadUint64At(data []byte, offset int) uint64 {
	if offset+32 > len(data) {
		return 0
	}
	return readUint64BE(data[offset+24 : offset+32])
}

func readUint64BE(b []byte) uint64 {
	return uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
}

func readUint32BE(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}
