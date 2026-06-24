package fastabi

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoder_Pool(t *testing.T) {
	e := NewEncoder()
	require.NotNil(t, e)
	PutEncoder(e)
	e2 := NewEncoder()
	require.NotNil(t, e2)
	PutEncoder(e2)
}

func TestEncoder_Reset(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint64(42)
	assert.NotEmpty(t, e.Bytes())
	e.Reset()
	assert.Empty(t, e.Bytes())
}

func TestEncoder_Hex(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint8(0xAB)
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000ab", e.Hex())
}

func TestEncoder_EncodeUint256_Nil(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint256(nil)
	assert.Equal(t, 32, len(e.Bytes()))
}

func TestEncoder_EncodeUint64(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint64(0xDEADBEEF)
	var exp [32]byte
	exp[28] = 0xDE
	exp[29] = 0xAD
	exp[30] = 0xBE
	exp[31] = 0xEF
	assert.Equal(t, exp[:], e.Bytes())
}

func TestEncoder_EncodeUint32(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint32(0xCAFEBABE)
	var exp [32]byte
	exp[28] = 0xCA
	exp[29] = 0xFE
	exp[30] = 0xBA
	exp[31] = 0xBE
	assert.Equal(t, exp[:], e.Bytes())
}

func TestEncoder_EncodeUint24(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint24(0xABCDEF)
	var exp [32]byte
	exp[29] = 0xAB
	exp[30] = 0xCD
	exp[31] = 0xEF
	assert.Equal(t, exp[:], e.Bytes())
}

func TestEncoder_EncodeUint16(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint16(0xABCD)
	var exp [32]byte
	exp[30] = 0xAB
	exp[31] = 0xCD
	assert.Equal(t, exp[:], e.Bytes())
}

func TestEncoder_EncodeUint8(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeUint8(0x42)
	var exp [32]byte
	exp[31] = 0x42
	assert.Equal(t, exp[:], e.Bytes())
}

func TestEncoder_EncodeAddress(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	var addr [20]byte
	addr[0] = 0xAA
	e.EncodeAddress(addr)
	assert.Equal(t, 32, len(e.Bytes()))
	assert.Equal(t, byte(0xAA), e.Bytes()[12])
}

func TestEncoder_EncodeBool(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeBool(true)
	assert.Equal(t, byte(1), e.Bytes()[31])
	e.Reset()
	e.EncodeBool(false)
	assert.Equal(t, byte(0), e.Bytes()[31])
}

func TestEncoder_EncodeBytes32(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	var b [32]byte
	b[0] = 0xAA
	e.EncodeBytes32(b)
	assert.Equal(t, b[:], e.Bytes())
}

func TestEncoder_EncodeMethodID(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeMethodID(MethodID{0x12, 0x34, 0x56, 0x78})
	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78}, e.Bytes())
}

func TestEncoder_EncodeBigInt_Nil(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeBigInt(nil)
	assert.Equal(t, 32, len(e.Bytes()))
}

func TestEncoder_EncodeBigInt_Positive(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeBigInt(big.NewInt(42))
	assert.Equal(t, byte(0x2A), e.Bytes()[31])
}

func TestEncoder_EncodeBigInt_Negative(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	e.EncodeBigInt(big.NewInt(-1))
	for _, b := range e.Bytes() {
		assert.Equal(t, byte(0xFF), b)
	}
}

func TestEncoder_EncodeBigInt_Large(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	val := new(big.Int).Lsh(big.NewInt(1), 200)
	e.EncodeBigInt(val)
	assert.Equal(t, 32, len(e.Bytes()))
}

func TestPutEncoder_LargeBuffer(t *testing.T) {
	e := NewEncoder()
	e.buf = make([]byte, 0, maxRetainCap+1)
	PutEncoder(e)
	e2 := NewEncoder()
	defer PutEncoder(e2)
	assert.LessOrEqual(t, cap(e2.buf), maxRetainCap)
}
