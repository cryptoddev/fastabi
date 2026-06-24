package fastabi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseType_Primitives(t *testing.T) {
	for _, tt := range []struct{ s string; k Kind }{
		{"uint256", KindUint256}, {"uint", KindUint256}, {"uint128", KindUint128},
		{"uint64", KindUint64}, {"uint32", KindUint32}, {"uint24", KindUint24},
		{"uint16", KindUint16}, {"uint8", KindUint8},
		{"int256", KindInt256}, {"int", KindInt256}, {"int128", KindInt128},
		{"int64", KindInt64}, {"int32", KindInt32}, {"int24", KindInt24},
		{"int16", KindInt16}, {"int8", KindInt8},
		{"address", KindAddress}, {"bool", KindBool}, {"bytes32", KindBytes32},
		{"bytes", KindBytes}, {"string", KindString},
	} {
		got, err := ParseType(tt.s)
		require.NoError(t, err, "ParseType(%q)", tt.s)
		assert.Equal(t, tt.k, got.Kind)
	}
}

func TestParseType_Arrays(t *testing.T) {
	t1, err := ParseType("uint256[]")
	require.NoError(t, err)
	assert.Equal(t, KindArray, t1.Kind)
	assert.Equal(t, KindUint256, t1.Elem.Kind)

	t2, err := ParseType("address[3]")
	require.NoError(t, err)
	assert.Equal(t, KindFixedArr, t2.Kind)
	assert.Equal(t, KindAddress, t2.Elem.Kind)
	assert.Equal(t, 3, t2.Size)
}

func TestParseType_Tuples(t *testing.T) {
	t1, err := ParseType("(uint8,address,bytes)")
	require.NoError(t, err)
	assert.Equal(t, KindTuple, t1.Kind)
	require.Len(t, t1.TupleEl, 3)

	t2, err := ParseType("(uint8,address,bytes)[]")
	require.NoError(t, err)
	assert.Equal(t, KindArray, t2.Kind)
	assert.Equal(t, KindTuple, t2.Elem.Kind)
}

func TestParseType_Errors(t *testing.T) {
	_, err := ParseType("")
	assert.Error(t, err)
	_, err = ParseType("unknown")
	assert.Error(t, err)
	_, err = ParseType("uint999")
	assert.Error(t, err)
	_, err = ParseType("int999")
	assert.Error(t, err)
	_, err = ParseType("bytes33")
	assert.Error(t, err)
	_, err = ParseType("bytes0")
	assert.Error(t, err)
	_, err = ParseType("uint7")
	assert.Error(t, err)
	_, err = ParseType("int7")
	assert.Error(t, err)
	_, err = ParseType("uint256[")
	assert.Error(t, err)
	_, err = ParseType("uint256[abc]")
	assert.Error(t, err)
}

func TestParseType_BytesN(t *testing.T) {
	pt, err := parseTypeStr("bytes1")
	require.NoError(t, err)
	assert.Equal(t, KindBytes32, pt.Kind)
	assert.Equal(t, 1, pt.Size)
}

func TestParseTypeStr_Empty(t *testing.T) {
	_, err := parseTypeStr("")
	assert.Error(t, err)
}

func TestParseTypes(t *testing.T) {
	ts, err := ParseTypes("uint256,address,bool")
	require.NoError(t, err)
	require.Len(t, ts, 3)
	assert.Equal(t, KindUint256, ts[0].Kind)
	assert.Equal(t, KindAddress, ts[1].Kind)
	assert.Equal(t, KindBool, ts[2].Kind)

	ts2, err := ParseTypes("")
	require.NoError(t, err)
	assert.Nil(t, ts2)

	ts3, err := ParseTypes("  ")
	require.NoError(t, err)
	assert.Nil(t, ts3)

	_, err = ParseTypes("uint256,unknown")
	assert.Error(t, err)
}

func TestParseFunction(t *testing.T) {
	name, inputs, err := ParseFunction("transfer(address,uint256)")
	require.NoError(t, err)
	assert.Equal(t, "transfer", name)
	require.Len(t, inputs, 2)
	assert.Equal(t, KindAddress, inputs[0].Kind)
	assert.Equal(t, KindUint256, inputs[1].Kind)

	name2, inputs2, err := ParseFunction("balanceOf(address):(uint256)")
	require.NoError(t, err)
	assert.Equal(t, "balanceOf", name2)
	require.Len(t, inputs2, 1)

	name3, inputs3, err := ParseFunction("myFunc")
	require.NoError(t, err)
	assert.Equal(t, "myFunc", name3)
	assert.Nil(t, inputs3)

	_, _, err = ParseFunction("myFunc(address")
	assert.Error(t, err)

	name4, inputs4, err := ParseFunction("execute((uint256,address),bool)")
	require.NoError(t, err)
	assert.Equal(t, "execute", name4)
	require.Len(t, inputs4, 2)
	assert.Equal(t, KindTuple, inputs4[0].Kind)
}

func TestParseTupleType_NoParen(t *testing.T) {
	_, err := parseTupleType("uint256")
	assert.Error(t, err)
}

func TestParseTupleType_UnmatchedParen(t *testing.T) {
	_, err := parseTupleType("(uint256,address")
	assert.Error(t, err)
}

func TestParseTupleType_InvalidInner(t *testing.T) {
	_, err := parseTupleType("(uint999)")
	assert.Error(t, err)
}

func TestParseTupleType_ArraySuffix(t *testing.T) {
	pt, err := parseTupleType("(uint256,address)[3]")
	require.NoError(t, err)
	assert.Equal(t, KindFixedArr, pt.Kind)
	assert.Equal(t, 3, pt.Size)
}

func TestParseTupleType_DynamicArraySuffix(t *testing.T) {
	pt, err := parseTupleType("(uint256)[]")
	require.NoError(t, err)
	assert.Equal(t, KindArray, pt.Kind)
}

func TestParseTupleType_UnmatchedBracket(t *testing.T) {
	_, err := parseTupleType("(uint256)[")
	assert.Error(t, err)
}

func TestParseTupleType_InvalidArraySize(t *testing.T) {
	_, err := parseTupleType("(uint256)[abc]")
	assert.Error(t, err)
}



func TestParamNames(t *testing.T) {
	ts := []ParamType{Named("from", TAddress()), Named("to", TAddress())}
	assert.Equal(t, []string{"from", "to"}, ParamNames(ts))

	ts2 := []ParamType{TAddress(), TUint256()}
	assert.Equal(t, []string{"address", "uint256"}, ParamNames(ts2))
}

func TestGoType(t *testing.T) {
	assert.Equal(t, "*U256", TUint256().GoType())
	assert.Equal(t, "*big.Int", TUint128().GoType())
	assert.Equal(t, "uint64", TUint64().GoType())
	assert.Equal(t, "int64", TInt64().GoType())
	assert.Equal(t, "[20]byte", TAddress().GoType())
	assert.Equal(t, "bool", TBool().GoType())
	assert.Equal(t, "[32]byte", TBytes32().GoType())
	assert.Equal(t, "[]byte", TBytes().GoType())
	assert.Equal(t, "string", TString().GoType())
	assert.Equal(t, "[]*U256", TArray(TUint256()).GoType())
	assert.Equal(t, "[]any", TTuple(TUint256()).GoType())
	assert.Equal(t, "any", ParamType{Kind: KindUnknown}.GoType())
}

func TestSignature(t *testing.T) {
	assert.Equal(t, "uint256", TUint256().Signature())
	assert.Equal(t, "address", TAddress().Signature())
	assert.Equal(t, "bool", TBool().Signature())
	assert.Equal(t, "bytes", TBytes().Signature())
	assert.Equal(t, "string", TString().Signature())
	assert.Equal(t, "uint256[]", TArray(TUint256()).Signature())
	assert.Equal(t, "uint256[3]", TFixedArray(TUint256(), 3).Signature())
	assert.Equal(t, "(uint256,address)", TTuple(TUint256(), TAddress()).Signature())
	assert.Equal(t, "unknown", ParamType{Kind: KindUnknown}.Signature())
}

func TestIsWhitespace(t *testing.T) {
	assert.True(t, IsWhitespace(" "))
	assert.True(t, IsWhitespace("\t\n\r"))
	assert.False(t, IsWhitespace("a"))
	assert.False(t, IsWhitespace("  x  "))
}
