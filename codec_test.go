package fastabi

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode_Uint256(t *testing.T) {
	val := NewU64(12345)
	encoded := Encode1(TUint256(), val)
	got, err := Decode1(TUint256(), encoded)
	require.NoError(t, err)
	assert.True(t, got.(*U256).Eq(val))
	PutU256(val)
	PutU256(got.(*U256))
}

func TestEncodeDecode_Address(t *testing.T) {
	expected := Address{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa,
		0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44}
	encoded := Encode1(TAddress(), expected)
	got, err := Decode1(TAddress(), encoded)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestEncodeDecode_Bool(t *testing.T) {
	for _, v := range []bool{true, false} {
		encoded := Encode1(TBool(), v)
		got, err := Decode1(TBool(), encoded)
		require.NoError(t, err)
		assert.Equal(t, v, got)
	}
}

func TestEncodeDecode_String(t *testing.T) {
	s := "Hello, World!"
	encoded := Encode1(TString(), s)
	got, err := Decode1(TString(), encoded)
	require.NoError(t, err)
	assert.Equal(t, s, got)
}

func TestEncodeDecode_Bytes(t *testing.T) {
	b := []byte{0x01, 0x02, 0x03}
	encoded := Encode1(TBytes(), b)
	got, err := Decode1(TBytes(), encoded)
	require.NoError(t, err)
	assert.Equal(t, b, got)
}

func TestEncodeDecode_EmptyBytes(t *testing.T) {
	encoded := Encode1(TBytes(), []byte{})
	got, err := Decode1(TBytes(), encoded)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestEncodeDecode_EmptyString(t *testing.T) {
	encoded := Encode1(TString(), "")
	got, err := Decode1(TString(), encoded)
	require.NoError(t, err)
	assert.Equal(t, "", got)
}

func TestEncodeDecode_Array(t *testing.T) {
	v1 := big.NewInt(100)
	v2 := big.NewInt(200)
	v3 := big.NewInt(300)
	vals := []any{v1, v2, v3}
	encoded := Encode1(TArray(TUint256()), vals)
	got, err := Decode1(TArray(TUint256()), encoded)
	require.NoError(t, err)
	gotSlice := got.([]any)
	require.Len(t, gotSlice, 3)
	assert.True(t, gotSlice[0].(*U256).Eq(NewU256FromBig(v1)))
	assert.True(t, gotSlice[1].(*U256).Eq(NewU256FromBig(v2)))
	assert.True(t, gotSlice[2].(*U256).Eq(NewU256FromBig(v3)))
	PutU256(gotSlice[0].(*U256))
	PutU256(gotSlice[1].(*U256))
	PutU256(gotSlice[2].(*U256))
}

func TestEncodeDecode_EmptyArray(t *testing.T) {
	encoded := Encode1(TArray(TUint256()), []any{})
	got, err := Decode1(TArray(TUint256()), encoded)
	require.NoError(t, err)
	assert.Empty(t, got.([]any))
}

func TestEncodeDecode_Tuple(t *testing.T) {
	addr := Address{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	val := NewU64(42)
	tuple := []any{uint8(7), addr, val}
	encoded := Encode1(TTuple(TUint8(), TAddress(), TUint256()), tuple)
	got, err := Decode1(TTuple(TUint8(), TAddress(), TUint256()), encoded)
	require.NoError(t, err)
	gotTuple := got.([]any)
	require.Len(t, gotTuple, 3)
	assert.Equal(t, uint64(7), gotTuple[0])
	assert.Equal(t, addr, gotTuple[1])
	assert.True(t, gotTuple[2].(*U256).Eq(val))
	PutU256(val)
	PutU256(gotTuple[2].(*U256))
}

func TestEncodeDecode_NestedTuple(t *testing.T) {
	inner := []any{NewU64(999), Address{0xaa}}
	outer := []any{inner, true}
	ts := TTuple(TTuple(TUint256(), TAddress()), TBool())
	encoded := Encode1(ts, outer)
	got, err := Decode1(ts, encoded)
	require.NoError(t, err)
	gotOuter := got.([]any)
	require.Len(t, gotOuter, 2)
	assert.Equal(t, true, gotOuter[1])
	gotInner := gotOuter[0].([]any)
	require.Len(t, gotInner, 2)
	assert.True(t, gotInner[0].(*U256).Eq(NewU64(999)))
	PutU256(gotInner[0].(*U256))
	PutU256(NewU64(999))
}

func TestEncodeDecode_ArrayOfTuples(t *testing.T) {
	tupleType := TTuple(TUint8(), TAddress(), TBytes())
	arrType := TArray(tupleType)
	tuple1 := []any{uint8(0), Address{1}, []byte{0x01, 0x02}}
	tuple2 := []any{uint8(1), Address{2}, []byte{0x03, 0x04, 0x05}}
	arr := []any{tuple1, tuple2}
	encoded := Encode1(arrType, arr)
	got, err := Decode1(arrType, encoded)
	require.NoError(t, err)
	gotArr := got.([]any)
	require.Len(t, gotArr, 2)
	gotT1 := gotArr[0].([]any)
	assert.Equal(t, uint64(0), gotT1[0])
	assert.Equal(t, tuple1[1], gotT1[1])
	assert.Equal(t, tuple1[2], gotT1[2])
}

func TestEncodeDecode_Multiple(t *testing.T) {
	expectedAddr := Address{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44}
	ts := []ParamType{TUint256(), TAddress(), TBool(), TString()}
	vals := []any{
		big.NewInt(12345),
		expectedAddr,
		true,
		"test string",
	}
	encoded := Encode(ts, vals)
	got, err := Decode(ts, encoded)
	require.NoError(t, err)
	require.Len(t, got, 4)
	assert.True(t, got[0].(*U256).Eq(NewU256FromBig(big.NewInt(12345))))
	assert.Equal(t, vals[1], got[1])
	assert.Equal(t, true, got[2])
	assert.Equal(t, "test string", got[3])
	PutU256(got[0].(*U256))
}

func TestEncodeDecode_Int256_Negative(t *testing.T) {
	v := big.NewInt(-1234567890)
	encoded := Encode1(TInt256(), v)
	got, err := Decode1(TInt256(), encoded)
	require.NoError(t, err)
	assert.Equal(t, v, got)
}

func TestEncodeDecode_Int24(t *testing.T) {
	for _, v := range []int64{-887272, -1000, -1, 0, 1, 1000, 887272} {
		encoded := Encode1(TInt24(), v)
		got, err := Decode1(TInt24(), encoded)
		require.NoError(t, err)
		assert.Equal(t, v, got.(int64))
	}
}

func TestEncodeDecode_FixedArray(t *testing.T) {
	vals := []any{big.NewInt(10), big.NewInt(20), big.NewInt(30)}
	encoded := Encode1(TFixedArray(TUint256(), 3), vals)
	got, err := Decode1(TFixedArray(TUint256(), 3), encoded)
	require.NoError(t, err)
	gotSlice := got.([]any)
	require.Len(t, gotSlice, 3)
	assert.True(t, gotSlice[0].(*U256).Eq(NewU256FromBig(big.NewInt(10))))
	assert.True(t, gotSlice[1].(*U256).Eq(NewU256FromBig(big.NewInt(20))))
	assert.True(t, gotSlice[2].(*U256).Eq(NewU256FromBig(big.NewInt(30))))
	PutU256(gotSlice[0].(*U256))
	PutU256(gotSlice[1].(*U256))
	PutU256(gotSlice[2].(*U256))
}

func TestEncodeDecode_Bytes32(t *testing.T) {
	b := [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}
	encoded := Encode1(TBytes32(), b)
	got, err := Decode1(TBytes32(), encoded)
	require.NoError(t, err)
	assert.Equal(t, b, got)
}

func TestEncodeDecode_Uint64(t *testing.T) {
	v := uint64(0xDEADBEEFCAFEBABE)
	encoded := Encode1(TUint64(), v)
	got, err := Decode1(TUint64(), encoded)
	require.NoError(t, err)
	assert.Equal(t, v, got)
}

func TestEncodeDecode_Uint8(t *testing.T) {
	v := uint8(42)
	encoded := Encode1(TUint8(), v)
	got, err := Decode1(TUint8(), encoded)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), got)
}

func TestEncodeDecode_Int128(t *testing.T) {
	v := new(big.Int)
	v.SetString("-12345678901234567890", 10)
	encoded := Encode1(TInt128(), v)
	got, err := Decode1(TInt128(), encoded)
	require.NoError(t, err)
	assert.Equal(t, v, got)
}

func TestEncodeDecode_Int64(t *testing.T) {
	v := int64(-9223372036854775807)
	encoded := Encode1(TInt64(), v)
	got, err := Decode1(TInt64(), encoded)
	require.NoError(t, err)
	assert.Equal(t, v, got)
}

func TestEncodeDecode_Int32(t *testing.T) {
	for _, v := range []int64{-2147483648, -1, 0, 1, 2147483647} {
		encoded := Encode1(TInt32(), v)
		got, err := Decode1(TInt32(), encoded)
		require.NoError(t, err)
		assert.Equal(t, v, got.(int64))
	}
}

func TestEncodeDecode_Int16(t *testing.T) {
	for _, v := range []int64{-32768, -1, 0, 1, 32767} {
		encoded := Encode1(TInt16(), v)
		got, err := Decode1(TInt16(), encoded)
		require.NoError(t, err)
		assert.Equal(t, v, got.(int64))
	}
}

func TestEncodeDecode_Int8(t *testing.T) {
	for _, v := range []int64{-128, -1, 0, 1, 127} {
		encoded := Encode1(TInt8(), v)
		got, err := Decode1(TInt8(), encoded)
		require.NoError(t, err)
		assert.Equal(t, v, got.(int64))
	}
}

func TestEncodeDecode_Uint128(t *testing.T) {
	v := new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)
	v.Lsh(v, 64)
	v.Add(v, new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF))
	encoded := Encode1(TUint128(), v)
	got, err := Decode1(TUint128(), encoded)
	require.NoError(t, err)
	assert.Equal(t, v, got)
}

func TestEncodeDecode_Uint32(t *testing.T) {
	v := uint64(0xFFFFFFFF)
	encoded := Encode1(TUint32(), v)
	got, err := Decode1(TUint32(), encoded)
	require.NoError(t, err)
	assert.Equal(t, v, got)
}

func TestEncodeDecode_Uint16(t *testing.T) {
	v := uint64(0xFFFF)
	encoded := Encode1(TUint16(), v)
	got, err := Decode1(TUint16(), encoded)
	require.NoError(t, err)
	assert.Equal(t, v, got)
}

func TestEncodeDecode_DirectArbPlan(t *testing.T) {
	actionType := TTuple(TUint8(), TAddress(), TBytes())
	planType := TTuple(TArray(actionType), TAddress(), TUint256(), TUint256(), TUint256())
	addr := Address{0xbb, 0x4C, 0xdB, 0x9C, 0xBd, 0x36, 0xB0, 0x1b, 0xD1, 0xcB,
		0xaE, 0xBF, 0x2D, 0xe0, 0x8d, 0x91, 0x73, 0xbc, 0x09, 0x5c}
	action := []any{uint8(0), Address{0x11}, []byte{0x01, 0x02}}
	plan := []any{[]any{action}, addr, big.NewInt(1e18), big.NewInt(0), big.NewInt(9999999999)}
	encoded := Encode1(planType, plan)
	got, err := Decode1(planType, encoded)
	require.NoError(t, err)
	gotPlan := got.([]any)
	require.Len(t, gotPlan, 5)
	assert.Equal(t, addr, gotPlan[1])
	PutU256(gotPlan[2].(*U256))
	PutU256(gotPlan[3].(*U256))
	PutU256(gotPlan[4].(*U256))
}

func TestEncodeDecode_Error(t *testing.T) {
	_, err := Decode1(TBytes(), []byte{0x01})
	assert.Error(t, err)

	_, err = Decode([]ParamType{TBytes()}, []byte{0x01})
	assert.Error(t, err)
}

func TestEncode_Empty(t *testing.T) {
	encoded := Encode(nil, nil)
	assert.Empty(t, encoded)
}

func TestEncodeStatic_DefaultKind(t *testing.T) {
	buf := encodeStatic(nil, ParamType{Kind: KindUnknown}, nil)
	assert.Equal(t, 32, len(buf))
}

func TestEncodeDynamic_DefaultKind(t *testing.T) {
	buf := encodeDynamic(nil, ParamType{Kind: KindUnknown}, nil)
	assert.Empty(t, buf)
}

func TestEncodeDynamic_String(t *testing.T) {
	buf := encodeDynamic(nil, TString(), "hello")
	assert.NotEmpty(t, buf)
}

func TestEncodeDynamic_Bytes(t *testing.T) {
	buf := encodeDynamic(nil, TBytes(), []byte{0x01})
	assert.NotEmpty(t, buf)
}

func TestEncodeDynamic_Array(t *testing.T) {
	vals := []any{big.NewInt(1)}
	buf := encodeDynamic(nil, TArray(TUint256()), vals)
	assert.NotEmpty(t, buf)
}

func TestEncodeDynamic_Tuple(t *testing.T) {
	val := []any{big.NewInt(1)}
	buf := encodeDynamic(nil, TTuple(TUint256()), val)
	assert.NotEmpty(t, buf)
}

func TestEncodeStatic_Tuple(t *testing.T) {
	val := []any{big.NewInt(1)}
	buf := encodeStatic(nil, TTuple(TUint256()), val)
	assert.Equal(t, 32, len(buf))
}

func TestEncodeStatic_FixedArray(t *testing.T) {
	vals := []any{big.NewInt(1)}
	buf := encodeStatic(nil, TFixedArray(TUint256(), 1), vals)
	assert.Equal(t, 32, len(buf))
}

func TestEncodeFixedArray_StaticElements(t *testing.T) {
	vals := []any{big.NewInt(1), big.NewInt(2)}
	buf := encodeFixedArray(nil, TFixedArray(TUint256(), 2), vals)
	assert.Equal(t, 64, len(buf))
}

func TestEncodeTupleVals(t *testing.T) {
	ts := []ParamType{TUint256(), TAddress()}
	vals := []any{big.NewInt(1), [20]byte{}}
	buf := encodeTupleVals(nil, ts, vals)
	assert.Equal(t, 64, len(buf))
}

func TestDynamicTailSize_Tuple(t *testing.T) {
	pt := TTuple(TUint256(), TAddress())
	val := []any{big.NewInt(1), [20]byte{}}
	assert.Equal(t, 64, dynamicTailSize(pt, val))
}

func TestAppendU256Val_LongBytes(t *testing.T) {
	bigVal := new(big.Int).Lsh(big.NewInt(1), 260)
	val := NewU256FromBig(bigVal)
	buf := appendU256Val(nil, val)
	assert.Equal(t, 32, len(buf))
	PutU256(val)
}

func TestAppendInt256Val_Negative(t *testing.T) {
	buf := appendInt256Val(nil, big.NewInt(-1))
	assert.Equal(t, 32, len(buf))
	for _, b := range buf {
		assert.Equal(t, byte(0xFF), b)
	}
}

func TestAppendInt256Val_Nil(t *testing.T) {
	buf := appendInt256Val(nil, nil)
	assert.Equal(t, 32, len(buf))
}

func TestToU256_Nil(t *testing.T) {
	r := toU256(nil)
	assert.True(t, r.IsZero())
	PutU256(r)
}

func TestToU256_Uint32(t *testing.T) {
	r := toU256(uint32(42))
	assert.True(t, r.Eq(NewU64(42)))
	PutU256(r)
}

func TestToU256_Uint16(t *testing.T) {
	r := toU256(uint16(42))
	assert.True(t, r.Eq(NewU64(42)))
	PutU256(r)
}

func TestToU256_Uint8(t *testing.T) {
	r := toU256(uint8(42))
	assert.True(t, r.Eq(NewU64(42)))
	PutU256(r)
}

func TestToU256_IntPositive(t *testing.T) {
	r := toU256(int(42))
	assert.True(t, r.Eq(NewU64(42)))
	PutU256(r)
}

func TestToU256_IntNegative(t *testing.T) {
	r := toU256(int(-1))
	assert.False(t, r.IsZero())
	PutU256(r)
}

func TestToU256_Default(t *testing.T) {
	r := toU256("42")
	assert.True(t, r.IsZero())
	PutU256(r)
}

func TestToBigInt_Nil(t *testing.T) {
	assert.Nil(t, toBigInt(nil))
}

func TestToBigInt_Uint64(t *testing.T) {
	r := toBigInt(uint64(42))
	assert.Equal(t, int64(42), r.Int64())
}

func TestToBigInt_Int(t *testing.T) {
	r := toBigInt(int(42))
	assert.Equal(t, int64(42), r.Int64())
}

func TestToBigInt_Default(t *testing.T) {
	r := toBigInt("42")
	assert.Equal(t, int64(0), r.Int64())
}

func TestToAddress_TypesAddress(t *testing.T) {
	addr := Address{0xBB}
	result := toAddress(addr)
	assert.Equal(t, byte(0xBB), result[0])
}

func TestToAddress_ShortBytes(t *testing.T) {
	result := toAddress([]byte{0xAA})
	assert.Equal(t, [20]byte{}, result)
}

func TestToBool_Uint64(t *testing.T) {
	assert.True(t, toBool(uint64(1)))
	assert.False(t, toBool(uint64(0)))
}

func TestToBool_Int(t *testing.T) {
	assert.True(t, toBool(int(1)))
	assert.False(t, toBool(int(0)))
}

func TestToBool_Default(t *testing.T) {
	assert.False(t, toBool("true"))
}

func TestToBytes32_Bytes(t *testing.T) {
	b := make([]byte, 32)
	b[0] = 0xAA
	result := toBytes32(b)
	assert.Equal(t, byte(0xAA), result[0])
}

func TestToBytes32_Default(t *testing.T) {
	result := toBytes32(42)
	assert.Equal(t, [32]byte{}, result)
}

func TestToBytes_String(t *testing.T) {
	result := toBytes("hello")
	assert.Equal(t, []byte("hello"), result)
}

func TestToBytes_Default(t *testing.T) {
	result := toBytes(42)
	assert.Nil(t, result)
}

func TestToString_Bytes(t *testing.T) {
	result := toString([]byte("hello"))
	assert.Equal(t, "hello", result)
}

func TestToString_Default(t *testing.T) {
	result := toString(42)
	assert.Equal(t, "", result)
}

func TestToSlice_StringSlice(t *testing.T) {
	result := toSlice([]string{"a", "b"})
	require.Len(t, result, 2)
	assert.Equal(t, "a", result[0])
}

func TestToSlice_BytesSlice(t *testing.T) {
	result := toSlice([][]byte{{0x01}})
	require.Len(t, result, 1)
}

func TestToSlice_AddressSlice(t *testing.T) {
	result := toSlice([]Address{{0xAA}})
	require.Len(t, result, 1)
}

func TestToSlice_U256Slice(t *testing.T) {
	result := toSlice([]*U256{NewU64(1)})
	require.Len(t, result, 1)
	PutU256(result[0].(*U256))
}

func TestToSlice_Default(t *testing.T) {
	assert.Nil(t, toSlice(42))
}
