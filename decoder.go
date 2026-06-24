package fastabi

import (
	"math/big"
	"sync"
)

const (
	defaultBufferCap = 1024
	maxRetainCap     = 4096
)

// Decoder is a reusable ABI decoder for hot-path static types.
// Uses zero-copy SetData() and sync.Pool for minimal allocations.
type Decoder struct {
	buf    []byte
	offset int
}

var decoderPool = sync.Pool{
	New: func() interface{} {
		return &Decoder{buf: make([]byte, 0, defaultBufferCap)}
	},
}

func GetDecoder() *Decoder {
	d := decoderPool.Get().(*Decoder)
	d.Reset()
	return d
}

func NewDecoder() *Decoder { return GetDecoder() }

func PutDecoder(d *Decoder) {
	d.offset = 0
	if cap(d.buf) > maxRetainCap {
		d.buf = make([]byte, 0, defaultBufferCap)
	} else {
		d.buf = d.buf[:0]
	}
	decoderPool.Put(d)
}

func (d *Decoder) Reset() {
	d.offset = 0
	d.buf = d.buf[:0]
}

// SetData uses zero-copy semantics. Caller must not mutate data while decoding.
func (d *Decoder) SetData(data []byte) {
	d.offset = 0
	d.buf = data
}

func (d *Decoder) Len() int { return len(d.buf) - d.offset }

func (d *Decoder) Offset() int { return d.offset }

func (d *Decoder) Skip(n int) bool {
	if n < 0 || d.offset+n > len(d.buf) {
		return false
	}
	d.offset += n
	return true
}

func (d *Decoder) remaining32() bool { return d.offset+32 <= len(d.buf) }

func (d *Decoder) DecodeUint256() *U256 {
	u := GetU256()
	if !d.remaining32() {
		return u
	}
	var b [32]byte
	copy(b[:], d.buf[d.offset:d.offset+32])
	d.offset += 32
	u.SetBytes32(b)
	return u
}

func (d *Decoder) DecodeAddress() [20]byte {
	var addr [20]byte
	if !d.remaining32() {
		return addr
	}
	copy(addr[:], d.buf[d.offset+12:d.offset+32])
	d.offset += 32
	return addr
}

func (d *Decoder) DecodeUint64() uint64 {
	if !d.remaining32() {
		return 0
	}
	b := d.buf[d.offset+24 : d.offset+32]
	d.offset += 32
	return uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
}

func (d *Decoder) DecodeUint32() uint32 {
	if !d.remaining32() {
		return 0
	}
	b := d.buf[d.offset+28 : d.offset+32]
	d.offset += 32
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func (d *Decoder) DecodeUint16() uint16 {
	if !d.remaining32() {
		return 0
	}
	b := d.buf[d.offset+30 : d.offset+32]
	d.offset += 32
	return uint16(b[0])<<8 | uint16(b[1])
}

func (d *Decoder) DecodeUint8() uint8 {
	if !d.remaining32() {
		return 0
	}
	v := d.buf[d.offset+31]
	d.offset += 32
	return v
}

func (d *Decoder) DecodeBool() bool {
	if !d.remaining32() {
		return false
	}
	v := d.buf[d.offset+31] != 0
	d.offset += 32
	return v
}

func (d *Decoder) DecodeBytes32() [32]byte {
	var out [32]byte
	if !d.remaining32() {
		return out
	}
	copy(out[:], d.buf[d.offset:d.offset+32])
	d.offset += 32
	return out
}

func (d *Decoder) DecodeBigInt() *big.Int {
	if !d.remaining32() {
		return big.NewInt(0)
	}
	v := new(big.Int).SetBytes(d.buf[d.offset : d.offset+32])
	d.offset += 32
	return v
}

// DecodeInt256 decodes signed int256 (two's complement).
func (d *Decoder) DecodeInt256() *big.Int {
	if !d.remaining32() {
		return big.NewInt(0)
	}
	var b [32]byte
	copy(b[:], d.buf[d.offset:d.offset+32])
	d.offset += 32
	if b[0]&0x80 == 0 {
		return new(big.Int).SetBytes(b[:])
	}
	for i := range b {
		b[i] = ^b[i]
	}
	v := new(big.Int).SetBytes(b[:])
	v.Add(v, big.NewInt(1))
	return v.Neg(v)
}
