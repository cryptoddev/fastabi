package fastabi

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecode_NilValues(t *testing.T) {
	tests := []struct {
		name string
		typ  ParamType
		val  any
	}{
		{"nil bytes", TBytes(), nil},
		{"nil string", TString(), nil},
		{"nil address", TAddress(), nil},
		{"nil bool", TBool(), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode([]ParamType{tt.typ}, []any{tt.val})
			require.NotNil(t, encoded)
			decoded, err := Decode([]ParamType{tt.typ}, encoded)
			require.NoError(t, err)
			require.Len(t, decoded, 1)
		})
	}
}

func TestEncodeDecode_EmptyBytesEdge(t *testing.T) {
	encoded := Encode([]ParamType{TBytes()}, []any{[]byte{}})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TBytes()}, encoded)
	require.NoError(t, err)
	assert.Empty(t, decoded[0].([]byte))
}

func TestEncodeDecode_EmptyStringEdge(t *testing.T) {
	encoded := Encode([]ParamType{TString()}, []any{""})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TString()}, encoded)
	require.NoError(t, err)
	assert.Equal(t, "", decoded[0].(string))
}

func TestEncodeDecode_U256ValueType(t *testing.T) {
	u := U256{}
	u.SetUint64(42)
	encoded := Encode([]ParamType{TUint256()}, []any{u})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TUint256()}, encoded)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), decoded[0].(*U256).Uint64())
}

func TestEncodeDecode_AddressFromString(t *testing.T) {
	addrHex := "0x1234567890abcdef1234567890abcdef12345678"
	addr, _ := ParseAddress(addrHex)
	encoded := Encode([]ParamType{TAddress()}, []any{addrHex})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TAddress()}, encoded)
	require.NoError(t, err)
	assert.Equal(t, addr, decoded[0].(Address))
}

func TestEncodeDecode_NegativeInt(t *testing.T) {
	v := big.NewInt(-42)
	encoded := Encode([]ParamType{TInt256()}, []any{v})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TInt256()}, encoded)
	require.NoError(t, err)
	assert.Equal(t, 0, decoded[0].(*big.Int).Cmp(big.NewInt(-42)))
}

func TestEncodeDecode_Int128Edge(t *testing.T) {
	v := big.NewInt(-1)
	v.Lsh(v, 127)
	encoded := Encode([]ParamType{TInt128()}, []any{v})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TInt128()}, encoded)
	require.NoError(t, err)

	expected := new(big.Int).Lsh(big.NewInt(-1), 127)
	assert.Equal(t, 0, decoded[0].(*big.Int).Cmp(expected))
}

func TestEncodeDecode_FixedArrayDynamicElems(t *testing.T) {
	typ := TFixedArray(TString(), 2)
	vals := []any{[]any{"hello", "world"}}
	encoded := Encode([]ParamType{typ}, vals)
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{typ}, encoded)
	require.NoError(t, err)
	arr := decoded[0].([]any)
	assert.Len(t, arr, 2)
	assert.Equal(t, "hello", arr[0].(string))
	assert.Equal(t, "world", arr[1].(string))
}

func TestEncodeDecode_NestedTuples(t *testing.T) {
	inner := TTuple(TAddress(), TBool())
	outer := TTuple(TUint256(), inner)
	addr, _ := ParseAddress("0xdead00000000000000000000000000000000beef")
	encoded := Encode([]ParamType{outer}, []any{[]any{
		NewU64(12345),
		[]any{addr, true},
	}})
	require.NotNil(t, encoded)

	decoded, err := Decode([]ParamType{outer}, encoded)
	require.NoError(t, err)
	require.Len(t, decoded, 1)
	outerVals := decoded[0].([]any)
	assert.Equal(t, uint64(12345), outerVals[0].(*U256).Uint64())
	innerVals := outerVals[1].([]any)
	assert.Equal(t, addr, innerVals[0].(Address))
	assert.True(t, innerVals[1].(bool))
}

func TestDecode_EmptyData(t *testing.T) {
	_, err := Decode([]ParamType{TUint256()}, nil)
	require.Error(t, err)

	_, err = Decode([]ParamType{TUint256()}, []byte{})
	require.Error(t, err)
}

func TestDecode_TruncatedData(t *testing.T) {
	encoded := Encode([]ParamType{TUint256()}, []any{NewU64(42)})
	_, err := Decode([]ParamType{TUint256()}, encoded[:16])
	require.Error(t, err)
}

func TestEncodeDecode_Int64Variants(t *testing.T) {
	tests := []struct {
		name string
		typ  ParamType
		val  int64
	}{
		{"int8", TInt8(), -128},
		{"int16", TInt16(), -32768},
		{"int32", TInt32(), -2147483648},
		{"int64", TInt64(), -9223372036854775808},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode([]ParamType{tt.typ}, []any{tt.val})
			require.NotNil(t, encoded)
			decoded, err := Decode([]ParamType{tt.typ}, encoded)
			require.NoError(t, err)
			assert.Equal(t, tt.val, decoded[0].(int64))
		})
	}
}

func TestEncodeDecode_Uint8Edge(t *testing.T) {
	for _, v := range []uint64{0, 1, 127, 128, 255} {
		encoded := Encode([]ParamType{TUint8()}, []any{v})
		require.NotNil(t, encoded)
		decoded, err := Decode([]ParamType{TUint8()}, encoded)
		require.NoError(t, err)
		assert.Equal(t, v, decoded[0].(uint64))
	}
}

func TestEncodeDecode_MaxUint256(t *testing.T) {
	maxVal := MaxU256()
	encoded := Encode([]ParamType{TUint256()}, []any{maxVal})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TUint256()}, encoded)
	require.NoError(t, err)
	assert.True(t, decoded[0].(*U256).Eq(maxVal))
}

func TestEncodeDecode_Bytes32Edge(t *testing.T) {
	var b32 [32]byte
	for i := range b32 {
		b32[i] = byte(i)
	}
	encoded := Encode([]ParamType{TBytes32()}, []any{b32})
	require.NotNil(t, encoded)
	decoded, err := Decode([]ParamType{TBytes32()}, encoded)
	require.NoError(t, err)
	assert.Equal(t, b32, decoded[0].([32]byte))
}
