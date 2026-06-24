package fastabi

import (
	"math/big"
)

func Encode(ts []ParamType, vals []any) []byte {
	if len(ts) != len(vals) {
		return nil
	}
	return encodeTuple(ts, vals)
}

func Encode1(t ParamType, val any) []byte {
	return encodeTuple([]ParamType{t}, []any{val})
}

func encodeTuple(ts []ParamType, vals []any) []byte {
	// First pass: compute head size and tail data
	headSize := 0
	for _, t := range ts {
		if t.IsDynamic() {
			headSize += 32
		} else {
			headSize += 32
		}
	}

	buf := make([]byte, 0, headSize+256)
	buf = encodeTupleHead(buf, ts, vals, headSize)
	buf = encodeTupleTail(buf, ts, vals)
	return buf
}

func encodeTupleHead(buf []byte, ts []ParamType, vals []any, headSize int) []byte {
	dynOffset := headSize
	for i, t := range ts {
		if t.IsDynamic() {
			buf = appendUint256FromUint64(buf, uint64(dynOffset))
			dynOffset += dynamicTailSize(t, vals[i])
		} else {
			buf = encodeStatic(buf, t, vals[i])
		}
	}
	return buf
}

func encodeTupleTail(buf []byte, ts []ParamType, vals []any) []byte {
	for i, t := range ts {
		if t.IsDynamic() {
			buf = encodeDynamic(buf, t, vals[i])
		}
	}
	return buf
}

func encodeTupleVals(buf []byte, ts []ParamType, vals []any) []byte {
	headSize := len(ts) * 32
	buf = encodeTupleHead(buf, ts, vals, headSize)
	buf = encodeTupleTail(buf, ts, vals)
	return buf
}

func encodeStatic(buf []byte, t ParamType, val any) []byte {
	switch t.Kind {
	case KindUint256, KindUint128, KindUint64, KindUint32, KindUint24, KindUint16, KindUint8:
		return appendU256Val(buf, toU256(val))
	case KindInt256, KindInt128, KindInt64, KindInt32, KindInt24, KindInt16, KindInt8:
		return appendInt256Val(buf, toBigInt(val))
	case KindAddress:
		return appendAddr(buf, toAddress(val))
	case KindBool:
		return appendBool(buf, toBool(val))
	case KindBytes32:
		return appendBytes32(buf, toBytes32(val))
	case KindFixedArr:
		return encodeFixedArray(buf, t, val)
	case KindTuple:
		return encodeTupleVals(buf, t.TupleEl, toSlice(val))
	default:
		return appendZeroes(buf, 32)
	}
}

func encodeDynamic(buf []byte, t ParamType, val any) []byte {
	switch t.Kind {
	case KindBytes:
		return appendDynBytes(buf, toBytes(val))
	case KindString:
		return appendDynBytes(buf, []byte(toString(val)))
	case KindArray:
		return encodeDynArray(buf, t, val)
	case KindFixedArr:
		return encodeFixedArray(buf, t, val)
	case KindTuple:
		return encodeTupleVals(buf, t.TupleEl, toSlice(val))
	default:
		return buf
	}
}

func encodeFixedArray(buf []byte, t ParamType, val any) []byte {
	slice := toSlice(val)
	if t.Elem.IsDynamic() {
		headSize := len(slice) * 32
		offsets := make([]int, len(slice))
		cum := headSize
		for i, v := range slice {
			offsets[i] = cum
			cum += dynamicTailSize(*t.Elem, v)
		}
		for _, off := range offsets {
			buf = appendUint256FromUint64(buf, uint64(off))
		}
		for _, v := range slice {
			buf = encodeDynamic(buf, *t.Elem, v)
		}
	} else {
		for _, v := range slice {
			buf = encodeStatic(buf, *t.Elem, v)
		}
	}
	return buf
}

func encodeDynArray(buf []byte, t ParamType, val any) []byte {
	slice := toSlice(val)
	buf = appendUint256FromUint64(buf, uint64(len(slice)))

	if t.Elem.IsDynamic() {
		headSize := len(slice) * 32
		offsets := make([]int, len(slice))
		cum := 32 + headSize
		for i, v := range slice {
			offsets[i] = cum
			cum += dynamicTailSize(*t.Elem, v)
		}
		for _, off := range offsets {
			buf = appendUint256FromUint64(buf, uint64(off))
		}
		for _, v := range slice {
			buf = encodeDynamic(buf, *t.Elem, v)
		}
	} else {
		for _, v := range slice {
			buf = encodeStatic(buf, *t.Elem, v)
		}
	}
	return buf
}

func dynamicTailSize(t ParamType, val any) int {
	if val == nil {
		return 32
	}
	switch t.Kind {
	case KindBytes:
		data := toBytes(val)
		if data == nil {
			return 32
		}
		return 32 + PaddedLen(len(data))
	case KindString:
		s := toString(val)
		return 32 + PaddedLen(len(s))
	case KindArray:
		return dynArrayTailSize(t, val)
	case KindFixedArr:
		return fixedArrayTailSize(t, val)
	case KindTuple:
		if t.TupleEl == nil {
			return 0
		}
		return tupleSize(t.TupleEl, toSlice(val))
	default:
		return 32
	}
}

func dynArrayTailSize(t ParamType, val any) int {
	slice := toSlice(val)
	size := 32
	if t.Elem != nil && t.Elem.IsDynamic() {
		headSize := len(slice) * 32
		for _, v := range slice {
			headSize += dynamicTailSize(*t.Elem, v)
		}
		size += headSize
	} else {
		size += len(slice) * 32
	}
	return size
}

func fixedArrayTailSize(t ParamType, val any) int {
	slice := toSlice(val)
	if t.Elem != nil && t.Elem.IsDynamic() {
		headSize := len(slice) * 32
		for _, v := range slice {
			headSize += dynamicTailSize(*t.Elem, v)
		}
		return headSize
	}
	return len(slice) * 32
}

func tupleSize(ts []ParamType, vals []any) int {
	size := 0
	for _, t := range ts {
		if t.IsDynamic() {
			size += 32
		} else {
			size += 32
		}
	}
	for i, t := range ts {
		if i < len(vals) && vals[i] != nil && t.IsDynamic() {
			size += dynamicTailSize(t, vals[i])
		}
	}
	return size
}

func PaddedLen(n int) int {
	return (n + 31) / 32 * 32
}

func appendU256Val(buf []byte, v *U256) []byte {
	if v == nil || v.IsZero() {
		return appendZeroes(buf, 32)
	}
	b := v.Bytes32()
	return append(buf, b[:]...)
}

func appendInt256Val(buf []byte, v *big.Int) []byte {
	if v == nil {
		return appendZeroes(buf, 32)
	}
	if v.Sign() >= 0 {
		b := v.Bytes()
		if len(b) > 32 {
			b = b[len(b)-32:]
		}
		pad := 32 - len(b)
		for i := 0; i < pad; i++ {
			buf = append(buf, 0)
		}
		return append(buf, b...)
	}
	// Negative: two's complement
	tmp := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), 256), v)
	b := tmp.Bytes()
	if len(b) > 32 {
		b = b[len(b)-32:]
	}
	pad := 32 - len(b)
	for i := 0; i < pad; i++ {
		buf = append(buf, 0)
	}
	return append(buf, b...)
}

func appendAddr(buf []byte, addr [20]byte) []byte {
	for i := 0; i < 12; i++ {
		buf = append(buf, 0)
	}
	return append(buf, addr[:]...)
}

func appendBool(buf []byte, v bool) []byte {
	for i := 0; i < 31; i++ {
		buf = append(buf, 0)
	}
	if v {
		return append(buf, 1)
	}
	return append(buf, 0)
}

func appendBytes32(buf []byte, b [32]byte) []byte {
	return append(buf, b[:]...)
}

func appendDynBytes(buf []byte, data []byte) []byte {
	buf = appendUint256FromUint64(buf, uint64(len(data)))
	buf = append(buf, data...)
	if rem := len(data) % 32; rem != 0 {
		for i := 0; i < 32-rem; i++ {
			buf = append(buf, 0)
		}
	}
	return buf
}

func appendZeroes(buf []byte, n int) []byte {
	for i := 0; i < n; i++ {
		buf = append(buf, 0)
	}
	return buf
}

func appendUint256FromUint64(buf []byte, v uint64) []byte {
	var slot [32]byte
	slot[24] = byte(v >> 56)
	slot[25] = byte(v >> 48)
	slot[26] = byte(v >> 40)
	slot[27] = byte(v >> 32)
	slot[28] = byte(v >> 24)
	slot[29] = byte(v >> 16)
	slot[30] = byte(v >> 8)
	slot[31] = byte(v)
	return append(buf, slot[:]...)
}

func toU256(v any) *U256 {
	switch x := v.(type) {
	case *U256:
		return x
	case U256:
		u := NewU64(0)
		u.Set(&x)
		return u
	case *big.Int:
		return NewU256FromBig(x)
	case uint64:
		return NewU64(x)
	case uint32:
		return NewU64(uint64(x))
	case uint16:
		return NewU64(uint64(x))
	case uint8:
		return NewU64(uint64(x))
	case int64:
		if x >= 0 {
			return NewU64(uint64(x))
		}
		return NewU256FromBig(big.NewInt(x))
	case int:
		if x >= 0 {
			return NewU64(uint64(x))
		}
		return NewU256FromBig(big.NewInt(int64(x)))
	default:
		return NewU64(0)
	}
}

func toBigInt(v any) *big.Int {
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case *big.Int:
		return x
	case *U256:
		return x.Inner().ToBig()
	case int64:
		return big.NewInt(x)
	case int:
		return big.NewInt(int64(x))
	case uint64:
		return new(big.Int).SetUint64(x)
	default:
		return big.NewInt(0)
	}
}

func toAddress(v any) [20]byte {
	switch x := v.(type) {
	case [20]byte:
		return x
	case Address:
		return [20]byte(x)
	case []byte:
		var a [20]byte
		if len(x) >= 20 {
			copy(a[:], x[len(x)-20:])
		}
		return a
	case string:
		addr, err := ParseAddress(x)
		if err == nil {
			return [20]byte(addr)
		}
		return [20]byte{}
	default:
		return [20]byte{}
	}
}

func toBool(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case uint64:
		return x != 0
	case int:
		return x != 0
	default:
		return false
	}
}

func toBytes32(v any) [32]byte {
	switch x := v.(type) {
	case [32]byte:
		return x
	case []byte:
		var b [32]byte
		copy(b[:], x)
		return b
	default:
		return [32]byte{}
	}
}

func toBytes(v any) []byte {
	switch x := v.(type) {
	case []byte:
		return x
	case string:
		return []byte(x)
	default:
		return nil
	}
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	default:
		return ""
	}
}

func toSlice(v any) []any {
	switch x := v.(type) {
	case []any:
		return x
	case []*U256:
		s := make([]any, len(x))
		for i, v := range x {
			s[i] = v
		}
		return s
	case []Address:
		s := make([]any, len(x))
		for i, v := range x {
			s[i] = v
		}
		return s
	case [][]byte:
		s := make([]any, len(x))
		for i, v := range x {
			s[i] = v
		}
		return s
	case []string:
		s := make([]any, len(x))
		for i, v := range x {
			s[i] = v
		}
		return s
	default:
		return nil
	}
}
