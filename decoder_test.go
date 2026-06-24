package fastabi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder_Pool(t *testing.T) {
	d := NewDecoder()
	require.NotNil(t, d)
	PutDecoder(d)
	d2 := NewDecoder()
	require.NotNil(t, d2)
	PutDecoder(d2)
}

func TestDecoder_SetData(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1, 2, 3, 4})
	assert.Equal(t, 4, d.Len())
	assert.Equal(t, 0, d.Offset())
}

func TestDecoder_Skip(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData(make([]byte, 64))
	assert.True(t, d.Skip(32))
	assert.Equal(t, 32, d.Offset())
	assert.False(t, d.Skip(-1))
	assert.False(t, d.Skip(33))
}

func TestDecoder_Reset(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData(make([]byte, 64))
	d.Skip(32)
	d.Reset()
	assert.Equal(t, 0, d.Offset())
	assert.Equal(t, 0, d.Len())
}

func TestDecoder_DecodeUint256(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[31] = 0x2A
	d.SetData(data)
	got := d.DecodeUint256()
	assert.True(t, got.Eq(NewU64(42)))
	PutU256(got)
}

func TestDecoder_DecodeUint256_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	got := d.DecodeUint256()
	assert.True(t, got.IsZero())
	PutU256(got)
}

func TestDecoder_DecodeAddress(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[12] = 0xAA
	d.SetData(data)
	got := d.DecodeAddress()
	assert.Equal(t, byte(0xAA), got[0])
}

func TestDecoder_DecodeAddress_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	got := d.DecodeAddress()
	assert.Equal(t, [20]byte{}, got)
}

func TestDecoder_DecodeUint64(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[31] = 0x2A
	d.SetData(data)
	assert.Equal(t, uint64(42), d.DecodeUint64())
}

func TestDecoder_DecodeUint64_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	assert.Equal(t, uint64(0), d.DecodeUint64())
}

func TestDecoder_DecodeUint32(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[31] = 0x2A
	d.SetData(data)
	assert.Equal(t, uint32(42), d.DecodeUint32())
}

func TestDecoder_DecodeUint32_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	assert.Equal(t, uint32(0), d.DecodeUint32())
}

func TestDecoder_DecodeUint16(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[31] = 0x2A
	d.SetData(data)
	assert.Equal(t, uint16(42), d.DecodeUint16())
}

func TestDecoder_DecodeUint16_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	assert.Equal(t, uint16(0), d.DecodeUint16())
}

func TestDecoder_DecodeUint8(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[31] = 0x2A
	d.SetData(data)
	assert.Equal(t, uint8(42), d.DecodeUint8())
}

func TestDecoder_DecodeUint8_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	assert.Equal(t, uint8(0), d.DecodeUint8())
}

func TestDecoder_DecodeBool(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	for _, v := range []bool{true, false} {
		data := make([]byte, 32)
		if v {
			data[31] = 1
		}
		d.SetData(data)
		assert.Equal(t, v, d.DecodeBool())
	}
}

func TestDecoder_DecodeBool_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	assert.False(t, d.DecodeBool())
}

func TestDecoder_DecodeBytes32(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[0] = 0xAA
	d.SetData(data)
	got := d.DecodeBytes32()
	assert.Equal(t, byte(0xAA), got[0])
}

func TestDecoder_DecodeBytes32_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	got := d.DecodeBytes32()
	assert.Equal(t, [32]byte{}, got)
}

func TestDecoder_DecodeBigInt(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[31] = 0x2A
	d.SetData(data)
	got := d.DecodeBigInt()
	assert.Equal(t, int64(42), got.Int64())
}

func TestDecoder_DecodeBigInt_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	got := d.DecodeBigInt()
	assert.Equal(t, int64(0), got.Int64())
}

func TestDecoder_DecodeInt256(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	data[31] = 0x2A
	d.SetData(data)
	got := d.DecodeInt256()
	assert.Equal(t, int64(42), got.Int64())
}

func TestDecoder_DecodeInt256_Negative(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	data := make([]byte, 32)
	for i := range data {
		data[i] = 0xFF
	}
	d.SetData(data)
	got := d.DecodeInt256()
	assert.Equal(t, int64(-1), got.Int64())
}

func TestDecoder_DecodeInt256_Overrun(t *testing.T) {
	d := NewDecoder()
	defer PutDecoder(d)
	d.SetData([]byte{1})
	got := d.DecodeInt256()
	assert.Equal(t, int64(0), got.Int64())
}

func TestPutDecoder_LargeBuffer(t *testing.T) {
	d := NewDecoder()
	d.buf = make([]byte, 0, maxRetainCap+1)
	PutDecoder(d)
	d2 := NewDecoder()
	defer PutDecoder(d2)
	assert.LessOrEqual(t, cap(d2.buf), maxRetainCap)
}
