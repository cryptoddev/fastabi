package fastabi

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeStatic_Address(t *testing.T) {
	data := make([]byte, 32)
	data[12] = 0xAA
	got, err := decodeStatic(TAddress(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, byte(0xAA), got.(Address)[0])
}

func TestDecodeStatic_Bool(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x01
	got, err := decodeStatic(TBool(), data, 0)
	require.NoError(t, err)
	assert.True(t, got.(bool))
}

func TestDecodeStatic_Bytes32(t *testing.T) {
	data := make([]byte, 32)
	data[0] = 0xAA
	got, err := decodeStatic(TBytes32(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, byte(0xAA), got.([32]byte)[0])
}

func TestDecodeStatic_Uint256(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TUint256(), data, 0)
	require.NoError(t, err)
	assert.True(t, got.(*U256).Eq(NewU64(42)))
	PutU256(got.(*U256))
}

func TestDecodeStatic_Uint128(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TUint128(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(*big.Int).Int64())
}

func TestDecodeStatic_Uint64(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TUint64(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), got)
}

func TestDecodeStatic_Uint32(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TUint32(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), got)
}

func TestDecodeStatic_Uint24(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TUint24(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), got)
}

func TestDecodeStatic_Uint16(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TUint16(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), got)
}

func TestDecodeStatic_Uint8(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TUint8(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), got)
}

func TestDecodeStatic_Int256(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TInt256(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(*big.Int).Int64())
}

func TestDecodeStatic_Int128(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TInt128(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(*big.Int).Int64())
}

func TestDecodeStatic_Int64(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TInt64(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(int64))
}

func TestDecodeStatic_Int32(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TInt32(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(int64))
}

func TestDecodeStatic_Int24(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TInt24(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(int64))
}

func TestDecodeStatic_Int16(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TInt16(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(int64))
}

func TestDecodeStatic_Int8(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TInt8(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(42), got.(int64))
}

func TestDecodeStatic_Tuple(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TTuple(TUint256()), data, 0)
	require.NoError(t, err)
	assert.True(t, got.([]any)[0].(*U256).Eq(NewU64(42)))
	PutU256(got.([]any)[0].(*U256))
}

func TestDecodeStatic_FixedArray(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeStatic(TFixedArray(TUint256(), 1), data, 0)
	require.NoError(t, err)
	assert.True(t, got.([]any)[0].(*U256).Eq(NewU64(42)))
	PutU256(got.([]any)[0].(*U256))
}

func TestDecodeStatic_Overrun(t *testing.T) {
	_, err := decodeStatic(TUint256(), []byte{0x01}, 0)
	assert.Error(t, err)
}

func TestDecodeStatic_DefaultKind(t *testing.T) {
	result, err := decodeStatic(ParamType{Kind: KindUnknown}, make([]byte, 32), 0)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestDecodeDynamic_Bytes(t *testing.T) {
	data := make([]byte, 64)
	data[31] = 2
	data[32] = 0x01
	data[33] = 0x02
	got, err := decodeDynamic(TBytes(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02}, got)
}

func TestDecodeDynamic_String(t *testing.T) {
	data := make([]byte, 64)
	data[31] = 5
	copy(data[32:37], "hello")
	got, err := decodeDynamic(TString(), data, 0)
	require.NoError(t, err)
	assert.Equal(t, "hello", got)
}

func TestDecodeDynamic_Array(t *testing.T) {
	data := make([]byte, 64)
	data[31] = 1
	data[63] = 0x2A
	got, err := decodeDynamic(TArray(TUint256()), data, 0)
	require.NoError(t, err)
	gotSlice := got.([]any)
	require.Len(t, gotSlice, 1)
	assert.True(t, gotSlice[0].(*U256).Eq(NewU64(42)))
	PutU256(gotSlice[0].(*U256))
}

func TestDecodeDynamic_Tuple(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeDynamic(TTuple(TUint256()), data, 0)
	require.NoError(t, err)
	gotSlice := got.([]any)
	assert.True(t, gotSlice[0].(*U256).Eq(NewU64(42)))
	PutU256(gotSlice[0].(*U256))
}

func TestDecodeDynamic_FixedArray(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := decodeDynamic(TFixedArray(TUint256(), 1), data, 0)
	require.NoError(t, err)
	gotSlice := got.([]any)
	require.Len(t, gotSlice, 1)
	assert.True(t, gotSlice[0].(*U256).Eq(NewU64(42)))
	PutU256(gotSlice[0].(*U256))
}

func TestDecodeDynamic_UnknownKind(t *testing.T) {
	_, err := decodeDynamic(ParamType{Kind: KindUnknown}, make([]byte, 32), 0)
	assert.Error(t, err)
}

func TestDecodeTuple_OffsetOverrun(t *testing.T) {
	_, err := decodeTuple([]ParamType{TBytes()}, []byte{0x01}, 0)
	assert.Error(t, err)
}

func TestDecodeTuple_StaticOverrun(t *testing.T) {
	_, err := decodeTuple([]ParamType{TUint256()}, []byte{0x01}, 0)
	assert.Error(t, err)
}

func TestDecodeTuple_DynamicNormal(t *testing.T) {
	// tuple (bytes): offset=32, length=2, data=0x01,0x02
	tupleEncoded := make([]byte, 96)
	tupleEncoded[31] = 32
	tupleEncoded[63] = 2
	tupleEncoded[64] = 0x01
	tupleEncoded[65] = 0x02
	got, err := decodeTuple([]ParamType{TBytes()}, tupleEncoded, 0)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02}, got[0].([]byte))
}

func TestDecodeDynBytes_LengthOverrun(t *testing.T) {
	assert.Nil(t, decodeDynBytes([]byte{0x01}, 0))
}

func TestDecodeDynBytes_DataOverrun(t *testing.T) {
	data := make([]byte, 64)
	data[31] = 100
	assert.Nil(t, decodeDynBytes(data, 0))
}

func TestDecodeDynArray_LengthOverrun(t *testing.T) {
	_, err := decodeDynArray(TArray(TUint256()), []byte{0x01}, 0)
	assert.Error(t, err)
}

func TestDecodeDynArray_OffsetsOverrun(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 1
	_, err := decodeDynArray(TArray(TBytes()), data, 0)
	assert.Error(t, err)
}

func TestDecodeFixedArray_StaticOverrun(t *testing.T) {
	data := make([]byte, 32)
	_, err := decodeFixedArray(TFixedArray(TUint256(), 2), data, 0)
	assert.Error(t, err)
}

func TestDecodeFixedArray_StaticNormal(t *testing.T) {
	data := make([]byte, 64)
	data[31] = 0x01
	data[63] = 0x02
	got, err := decodeFixedArray(TFixedArray(TUint256(), 2), data, 0)
	require.NoError(t, err)
	require.Len(t, got, 2)
	PutU256(got[0].(*U256))
	PutU256(got[1].(*U256))
}

func TestDecodeInt256_Positive(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	assert.Equal(t, int64(42), decodeInt256(data).Int64())
}

func TestDecodeInt256_Negative(t *testing.T) {
	data := make([]byte, 32)
	for i := range data {
		data[i] = 0xFF
	}
	assert.Equal(t, int64(-1), decodeInt256(data).Int64())
}

func TestDecodeInt128_Positive(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	assert.Equal(t, int64(42), decodeInt128(data).Int64())
}

func TestDecodeInt128_Negative(t *testing.T) {
	data := make([]byte, 32)
	for i := 16; i < 32; i++ {
		data[i] = 0xFF
	}
	assert.Equal(t, int64(-1), decodeInt128(data).Int64())
}

func TestReadUint64At_ShortData(t *testing.T) {
	assert.Equal(t, uint64(0), ReadUint64At([]byte{0x01}, 0))
}

func TestReadUint64At_Normal(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	assert.Equal(t, uint64(42), ReadUint64At(data, 0))
}

func TestReadUint32BE(t *testing.T) {
	b := []byte{0x00, 0x00, 0x00, 0x2A}
	assert.Equal(t, uint32(42), readUint32BE(b))
}
