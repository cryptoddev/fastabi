package fastabi

import (
	"encoding/hex"
	"math/big"
	"sync"
)

type Encoder struct {
	buf []byte
}

var encoderPool = sync.Pool{
	New: func() interface{} {
		return &Encoder{buf: make([]byte, 0, defaultBufferCap)}
	},
}

func GetEncoder() *Encoder {
	e := encoderPool.Get().(*Encoder)
	e.Reset()
	return e
}

func NewEncoder() *Encoder { return GetEncoder() }

func PutEncoder(e *Encoder) {
	if cap(e.buf) > maxRetainCap {
		e.buf = make([]byte, 0, defaultBufferCap)
	} else {
		e.buf = e.buf[:0]
	}
	encoderPool.Put(e)
}

func (e *Encoder) Reset() { e.buf = e.buf[:0] }

func (e *Encoder) Bytes() []byte { return e.buf }

func (e *Encoder) Hex() string { return "0x" + hex.EncodeToString(e.buf) }

func (e *Encoder) appendZeroes(n int) {
	for i := 0; i < n; i++ {
		e.buf = append(e.buf, 0)
	}
}

func (e *Encoder) EncodeUint256(u *U256) *Encoder {
	if u == nil {
		e.appendZeroes(32)
		return e
	}
	b := u.Bytes()
	if len(b) > 32 {
		b = b[len(b)-32:]
	}
	if pad := 32 - len(b); pad > 0 {
		e.appendZeroes(pad)
	}
	e.buf = append(e.buf, b...)
	return e
}

func (e *Encoder) EncodeUint64(v uint64) *Encoder {
	e.appendZeroes(24)
	e.buf = append(e.buf, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	return e
}

func (e *Encoder) EncodeUint32(v uint32) *Encoder {
	e.appendZeroes(28)
	e.buf = append(e.buf, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	return e
}

func (e *Encoder) EncodeUint24(v uint32) *Encoder {
	e.appendZeroes(29)
	e.buf = append(e.buf, byte(v>>16), byte(v>>8), byte(v))
	return e
}

func (e *Encoder) EncodeUint16(v uint16) *Encoder {
	e.appendZeroes(30)
	e.buf = append(e.buf, byte(v>>8), byte(v))
	return e
}

func (e *Encoder) EncodeUint8(v uint8) *Encoder {
	e.appendZeroes(31)
	e.buf = append(e.buf, v)
	return e
}

func (e *Encoder) EncodeAddress(addr [20]byte) *Encoder {
	e.appendZeroes(12)
	e.buf = append(e.buf, addr[:]...)
	return e
}

func (e *Encoder) EncodeBool(v bool) *Encoder {
	e.appendZeroes(31)
	if v {
		e.buf = append(e.buf, 1)
	} else {
		e.buf = append(e.buf, 0)
	}
	return e
}

func (e *Encoder) EncodeBytes32(v [32]byte) *Encoder {
	e.buf = append(e.buf, v[:]...)
	return e
}

type MethodID [4]byte

func (e *Encoder) EncodeMethodID(id MethodID) *Encoder {
	e.buf = append(e.buf, id[:]...)
	return e
}

var twoTo256 = new(big.Int).Lsh(big.NewInt(1), 256)

func (e *Encoder) EncodeBigInt(v *big.Int) *Encoder {
	if v == nil {
		e.appendZeroes(32)
		return e
	}
	if v.Sign() >= 0 {
		b := v.Bytes()
		if len(b) > 32 {
			b = b[len(b)-32:]
		}
		if pad := 32 - len(b); pad > 0 {
			e.appendZeroes(pad)
		}
		e.buf = append(e.buf, b...)
		return e
	}
	// negative => two's complement using scratch buffer
	tmp := new(big.Int).Add(twoTo256, v)
	b := tmp.Bytes()
	if len(b) > 32 {
		b = b[len(b)-32:]
	}
	if pad := 32 - len(b); pad > 0 {
		e.appendZeroes(pad)
	}
	e.buf = append(e.buf, b...)
	return e
}
