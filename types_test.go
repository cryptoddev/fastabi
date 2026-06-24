package fastabi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKind_String(t *testing.T) {
	assert.Equal(t, "uint8", KindUint8.String())
	assert.Equal(t, "uint16", KindUint16.String())
	assert.Equal(t, "uint24", KindUint24.String())
	assert.Equal(t, "uint32", KindUint32.String())
	assert.Equal(t, "uint64", KindUint64.String())
	assert.Equal(t, "uint128", KindUint128.String())
	assert.Equal(t, "uint256", KindUint256.String())
	assert.Equal(t, "int8", KindInt8.String())
	assert.Equal(t, "int16", KindInt16.String())
	assert.Equal(t, "int24", KindInt24.String())
	assert.Equal(t, "int32", KindInt32.String())
	assert.Equal(t, "int64", KindInt64.String())
	assert.Equal(t, "int128", KindInt128.String())
	assert.Equal(t, "int256", KindInt256.String())
	assert.Equal(t, "address", KindAddress.String())
	assert.Equal(t, "bool", KindBool.String())
	assert.Equal(t, "bytes32", KindBytes32.String())
	assert.Equal(t, "bytes", KindBytes.String())
	assert.Equal(t, "string", KindString.String())
	assert.Equal(t, "array", KindArray.String())
	assert.Equal(t, "fixedArray", KindFixedArr.String())
	assert.Equal(t, "tuple", KindTuple.String())
	assert.Equal(t, "unknown", KindUnknown.String())
}

func TestIsDynamic(t *testing.T) {
	assert.False(t, TUint256().IsDynamic())
	assert.False(t, TAddress().IsDynamic())
	assert.False(t, TBool().IsDynamic())
	assert.True(t, TBytes().IsDynamic())
	assert.True(t, TString().IsDynamic())
	assert.True(t, TArray(TUint256()).IsDynamic())
	assert.False(t, TFixedArray(TUint256(), 3).IsDynamic())
	assert.True(t, TFixedArray(TBytes(), 3).IsDynamic())
	assert.False(t, TTuple(TUint256(), TAddress()).IsDynamic())
	assert.True(t, TTuple(TUint256(), TBytes()).IsDynamic())
	assert.False(t, ParamType{Kind: KindTuple, TupleEl: nil}.IsDynamic())
}

func TestStaticSize(t *testing.T) {
	assert.Equal(t, 32, TUint256().StaticSize())
	assert.Equal(t, 96, TTuple(TUint256(), TAddress(), TBool()).StaticSize())
	assert.Equal(t, 96, TFixedArray(TUint256(), 3).StaticSize())
	assert.Equal(t, 0, TBytes().StaticSize())
	assert.Equal(t, 0, ParamType{Kind: KindFixedArr, Elem: nil, Size: 3}.StaticSize())
}

func TestAddress(t *testing.T) {
	addr := Address{0x12, 0x34, 0x56, 0x78}

	str := addr.String()
	if str == "" {
		t.Error("Address.String returned empty")
	}

	zero := Address{}
	if !zero.IsZero() {
		t.Error("Zero address should be zero")
	}

	if addr.IsZero() {
		t.Error("Non-zero address should not be zero")
	}
}

func TestHash(t *testing.T) {
	h, err := ParseHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Fatalf("ParseHash failed: %v", err)
	}
	if !strings.HasPrefix(h.Hex(), "0x") {
		t.Error("Hash.Hex should start with 0x")
	}
	if h.IsZero() {
		t.Error("Non-zero hash should not be zero")
	}
	if len(h.Bytes()) != 32 {
		t.Errorf("Hash.Bytes length: expected 32, got %d", len(h.Bytes()))
	}

	zero := Hash{}
	if !zero.IsZero() {
		t.Error("Zero hash should be zero")
	}
}

func TestParseHash_Errors(t *testing.T) {
	// Bad length
	_, err := ParseHash("0x1234")
	if err == nil {
		t.Error("expected error for short hash")
	}

	// Invalid hex
	_, err = ParseHash("0xGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG")
	if err == nil {
		t.Error("expected error for invalid hex")
	}

	// Valid no-prefix
	_, err = ParseHash("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	if err != nil {
		t.Errorf("unexpected error for valid hash: %v", err)
	}
}
