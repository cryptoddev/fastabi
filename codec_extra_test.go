package fastabi

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode_Error(t *testing.T) {
	// mismatched lengths
	assert.Nil(t, Encode([]ParamType{TUint256()}, nil))
}

func TestEncodeFixedArray_DynamicElements(t *testing.T) {
	vals := []any{[]byte{0x01}, []byte{0x02, 0x03}}
	buf := encodeFixedArray(nil, TFixedArray(TBytes(), 2), vals)
	assert.NotEmpty(t, buf)
}

func TestAppendU256Val_Nil(t *testing.T) {
	buf := appendU256Val(nil, nil)
	assert.Equal(t, 32, len(buf))
}

func TestAppendInt256Val_Positive(t *testing.T) {
	buf := appendInt256Val(nil, big.NewInt(42))
	assert.Equal(t, 32, len(buf))
	assert.Equal(t, byte(0x2A), buf[31])
}

func TestToU256_BigInt(t *testing.T) {
	r := toU256(big.NewInt(42))
	assert.True(t, r.Eq(NewU64(42)))
	PutU256(r)
}

func TestToU256_Int64Positive(t *testing.T) {
	r := toU256(int64(42))
	assert.True(t, r.Eq(NewU64(42)))
	PutU256(r)
}

func TestToU256_Int64Negative(t *testing.T) {
	r := toU256(int64(-1))
	assert.False(t, r.IsZero())
	PutU256(r)
}

func TestToAddress_NilDefault(t *testing.T) {
	result := toAddress(nil)
	assert.Equal(t, [20]byte{}, result)
}

func TestToBool_NilDefault(t *testing.T) {
	assert.False(t, toBool(nil))
}

func TestToBytes32_NilDefault(t *testing.T) {
	result := toBytes32(nil)
	assert.Equal(t, [32]byte{}, result)
}

func TestToString_NilDefault(t *testing.T) {
	result := toString(nil)
	assert.Equal(t, "", result)
}

func TestToSlice_NilDefault(t *testing.T) {
	assert.Nil(t, toSlice(nil))
}

func TestToBigInt_U256(t *testing.T) {
	u := NewU64(42)
	r := toBigInt(u)
	assert.Equal(t, int64(42), r.Int64())
	PutU256(u)
}

func TestEncode1(t *testing.T) {
	encoded := Encode1(TUint256(), NewU64(1))
	assert.Equal(t, 32, len(encoded))
	PutU256(NewU64(1))
}

func TestDecode1(t *testing.T) {
	data := make([]byte, 32)
	data[31] = 0x2A
	got, err := Decode1(TUint256(), data)
	assert.NoError(t, err)
	assert.True(t, got.(*U256).Eq(NewU64(42)))
	PutU256(got.(*U256))
}

func TestDecode1_Error(t *testing.T) {
	_, err := Decode1(TTuple(TBytes(), TAddress()), make([]byte, 32))
	assert.Error(t, err)
}

func TestDecode_Error(t *testing.T) {
	_, err := Decode([]ParamType{TBytes()}, []byte{0x01})
	assert.Error(t, err)
}

func TestSignature_FixedArray(t *testing.T) {
	assert.Equal(t, "uint256[3]", TFixedArray(TUint256(), 3).Signature())
}

func TestSignature_DynamicArray(t *testing.T) {
	assert.Equal(t, "bytes[]", TArray(TBytes()).Signature())
}

func TestSignature_Tuple(t *testing.T) {
	assert.Equal(t, "(uint256,address)", TTuple(TUint256(), TAddress()).Signature())
}

func TestSignature_ArrayOfTuple(t *testing.T) {
	assert.Equal(t, "(uint256)[]", TArray(TTuple(TUint256())).Signature())
}

func TestSignature_BytesN(t *testing.T) {
	assert.Equal(t, "bytes1", ParamType{Kind: KindBytes32, Size: 1}.Signature())
}

func TestSignature_ArrayNilElem(t *testing.T) {
	assert.Equal(t, "any[]", ParamType{Kind: KindArray}.Signature())
}

func TestSignature_FixedArrayNilElem(t *testing.T) {
	assert.Equal(t, "any[]", ParamType{Kind: KindFixedArr, Size: 3}.Signature())
}

func TestGoType_FixedArray(t *testing.T) {
	assert.Equal(t, "[3]*U256", TFixedArray(TUint256(), 3).GoType())
}

func TestGoType_ArrayNilElem(t *testing.T) {
	assert.Equal(t, "[]any", ParamType{Kind: KindArray}.GoType())
}

func TestGoType_FixedArrayNilElem(t *testing.T) {
	assert.Equal(t, "[]any", ParamType{Kind: KindFixedArr, Size: 3}.GoType())
}


func TestParseType_UintN(t *testing.T) {
	pt, err := ParseType("uint128")
	assert.NoError(t, err)
	assert.Equal(t, KindUint128, pt.Kind)
}

func TestParseType_IntN(t *testing.T) {
	pt, err := ParseType("int128")
	assert.NoError(t, err)
	assert.Equal(t, KindInt128, pt.Kind)
}

func TestEncoder_EncodeUint256_NonNil(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	val := NewU64(42)
	e.EncodeUint256(val)
	assert.Equal(t, byte(0x2A), e.Bytes()[31])
	PutU256(val)
}

func TestEncoder_EncodeBigInt_Overflow(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	// Value > 256 bits gets truncated
	bigVal := new(big.Int).Lsh(big.NewInt(1), 260)
	e.EncodeBigInt(bigVal)
	assert.Equal(t, 32, len(e.Bytes()))
}


func TestDecodeDynBytes_Normal(t *testing.T) {
	data := make([]byte, 64)
	data[31] = 2
	data[32] = 0xAB
	data[33] = 0xCD
	result := decodeDynBytes(data, 0)
	assert.Equal(t, []byte{0xAB, 0xCD}, result)
}

func TestDecodeDynArray_StaticElements(t *testing.T) {
	data := make([]byte, 96)
	data[31] = 2       // length = 2
	data[63] = 0x01    // first element value
	data[95] = 0x2A    // second element value
	got, err := decodeDynArray(TArray(TUint256()), data, 0)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	PutU256(got[0].(*U256))
	PutU256(got[1].(*U256))
}


func TestDynamicTailSize_Bytes(t *testing.T) {
	size := dynamicTailSize(TBytes(), []byte{0x01, 0x02})
	assert.Equal(t, 64, size) // 32 (length) + 32 (padded data)
}

func TestDynamicTailSize_String(t *testing.T) {
	size := dynamicTailSize(TString(), "hi")
	assert.Equal(t, 64, size)
}

func TestDynamicTailSize_Array(t *testing.T) {
	size := dynamicTailSize(TArray(TUint256()), []any{big.NewInt(1), big.NewInt(2)})
	assert.Equal(t, 96, size) // 32 (length) + 64 (2 * 32)
}

func TestDynamicTailSize_Default(t *testing.T) {
	size := dynamicTailSize(TAddress(), [20]byte{})
	assert.Equal(t, 32, size)
}

func TestTupleSize(t *testing.T) {
	ts := []ParamType{TUint256(), TAddress()}
	vals := []any{big.NewInt(1), [20]byte{}}
	assert.Equal(t, 64, tupleSize(ts, vals))
}

func TestTupleSize_Dynamic(t *testing.T) {
	ts := []ParamType{TBytes()}
	vals := []any{[]byte{0x01}}
	assert.Equal(t, 96, tupleSize(ts, vals)) // 32 offset + 32 length + 32 padded data
}

func TestAppendU256Val_Zero(t *testing.T) {
	u := NewU64(0)
	buf := appendU256Val(nil, u)
	assert.Equal(t, 32, len(buf))
	PutU256(u)
}

func TestToAddress_BytesLong(t *testing.T) {
	b := make([]byte, 25)
	b[5] = 0xAA
	result := toAddress(b)
	assert.Equal(t, byte(0xAA), result[0])
}

func TestSignature_BytesN_Zero(t *testing.T) {
	assert.Equal(t, "bytes32", TBytes32().Signature())
}

func TestDecodeTuple_DynamicErrorPropagation(t *testing.T) {
	data2 := make([]byte, 64)
	data2[31] = 32
	data2[63] = 1 // array length = 1, but no element data
	_, err := decodeTuple([]ParamType{TArray(TUint256())}, data2, 0)
	assert.Error(t, err)
}

func TestDecodeDynArray_StaticOverrun(t *testing.T) {
	data := make([]byte, 40)
	data[31] = 2 // length = 2, but only 8 bytes of element data
	_, err := decodeDynArray(TArray(TUint256()), data, 0)
	assert.Error(t, err)
}

func TestEncoder_EncodeUint256_LongBytes(t *testing.T) {
	e := NewEncoder()
	defer PutEncoder(e)
	bigVal := new(big.Int).Lsh(big.NewInt(1), 260)
	val := NewU256FromBig(bigVal)
	e.EncodeUint256(val)
	assert.Equal(t, 32, len(e.Bytes()))
	PutU256(val)
}

func TestAppendInt256Val_LongBytes(t *testing.T) {
	bigVal := new(big.Int).Lsh(big.NewInt(1), 260)
	buf := appendInt256Val(nil, bigVal)
	assert.Equal(t, 32, len(buf))
}

func TestParseType_UintN_Generic(t *testing.T) {
	pt, err := parseTypeStr("uint128")
	assert.NoError(t, err)
	assert.Equal(t, KindUint128, pt.Kind)
}

func TestParseType_IntN_Generic(t *testing.T) {
	pt, err := parseTypeStr("int128")
	assert.NoError(t, err)
	assert.Equal(t, KindInt128, pt.Kind)
}

func TestMustDecodeSignature_Panic(t *testing.T) {
	defer func() { _ = recover() }()
	MustDecodeSignature("test()")
}
