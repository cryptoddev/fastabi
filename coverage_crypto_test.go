package fastabi

import (
	"math/big"
	"testing"
)

func TestU256_SubUnderflow(t *testing.T) {
	a := NewU64(100)
	b := NewU64(200)
	a.Sub(b)
	_ = a.Uint64()
}

func TestU256_MulDiv_Cover(t *testing.T) {
	a := NewU64(1e18)
	b := NewU64(2e18)
	c := NewU64(1e18)
	result := GetU256().Set(a)
	result.MulDiv(b, c)
	_ = result.Uint64()
	PutU256(result)
	PutU256(a)
	PutU256(b)
	PutU256(c)
}

func TestU256_MulDivOverflow_Cover(t *testing.T) {
	a := NewU256FromBig(new(big.Int).Lsh(big.NewInt(1), 200))
	b := NewU256FromBig(new(big.Int).Lsh(big.NewInt(1), 200))
	c := NewU256FromBig(new(big.Int).Lsh(big.NewInt(1), 100))
	result := GetU256().Set(a)
	result.MulDiv(b, c)
	PutU256(result)
	PutU256(a)
	PutU256(b)
	PutU256(c)
}

func TestU256_MulDivRoundingUp_Cover(t *testing.T) {
	a := NewU64(10)
	b := NewU64(7)
	c := NewU64(3)
	result := GetU256().Set(a)
	result.MulDivRoundingUp(b, c)
	PutU256(result)
	PutU256(a)
	PutU256(b)
	PutU256(c)
}

func TestU256_MulDivRoundingUpExact_Cover(t *testing.T) {
	a := NewU64(10)
	b := NewU64(5)
	c := NewU64(2)
	result := GetU256().Set(a)
	result.MulDivRoundingUp(b, c)
	PutU256(result)
	PutU256(a)
	PutU256(b)
	PutU256(c)
}

func TestU256_MulDivRoundingUpOverflow_Cover(t *testing.T) {
	a := NewU256FromBig(new(big.Int).Lsh(big.NewInt(1), 200))
	b := NewU256FromBig(new(big.Int).Lsh(big.NewInt(1), 200))
	c := NewU256FromBig(new(big.Int).Lsh(big.NewInt(1), 100))
	result := GetU256().Set(a)
	result.MulDivRoundingUp(b, c)
	PutU256(result)
	PutU256(a)
	PutU256(b)
	PutU256(c)
}

func TestHALF_WAD_RAYInt(t *testing.T) {
	_ = HALF_WAD()
	_ = RAYInt()
}

func TestAddressFromUint64(t *testing.T) {
	addr := AddressFromUint64(0xDEADBEEF)
	if addr.IsZero() {
		t.Fatal("expected non-zero address")
	}
}

func TestParseHexByte_Invalid(t *testing.T) {
	_, err := ParseHexByte('z', 'z')
	if err == nil {
		t.Fatal("expected error for invalid hex")
	}
	_, err = ParseHexByte('0', 'z')
	if err == nil {
		t.Fatal("expected error for invalid hex")
	}
}

func TestParseHexUint64_Invalid(t *testing.T) {
	_, err := ParseHexUint64("0xzzzz")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHexToBytes_OddLength(t *testing.T) {
	b, err := HexToBytes("0xa")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 1 || b[0] != 0x0a {
		t.Fatal("expected 0x0a")
	}
}

func TestSignature_AllTypes(t *testing.T) {
	tests := []struct {
		name string
		typ  ParamType
		want string
	}{
		{"uint8", TUint8(), "uint8"},
		{"uint16", TUint16(), "uint16"},
		{"uint24", TUint24(), "uint24"},
		{"uint32", TUint32(), "uint32"},
		{"uint64", TUint64(), "uint64"},
		{"uint128", TUint128(), "uint128"},
		{"uint256", TUint256(), "uint256"},
		{"int8", TInt8(), "int8"},
		{"int16", TInt16(), "int16"},
		{"int24", TInt24(), "int24"},
		{"int32", TInt32(), "int32"},
		{"int64", TInt64(), "int64"},
		{"int128", TInt128(), "int128"},
		{"int256", TInt256(), "int256"},
		{"address", TAddress(), "address"},
		{"bool", TBool(), "bool"},
		{"bytes", TBytes(), "bytes"},
		{"string", TString(), "string"},
		{"bytes32", TBytes32(), "bytes32"},
		{"bytes8", ParamType{Kind: KindBytes32, Size: 8}, "bytes8"},
		{"uint256[]", TArray(TUint256()), "uint256[]"},
		{"uint256[3]", TFixedArray(TUint256(), 3), "uint256[3]"},
		{"tuple(uint256,address)", TTuple(TUint256(), TAddress()), "(uint256,address)"},
		{"tuple(uint256,address)[]", TArray(TTuple(TUint256(), TAddress())), "(uint256,address)[]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.typ.Signature()
			if got != tt.want {
				t.Errorf("Signature() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSignature_Unknown(t *testing.T) {
	var p ParamType
	got := p.Signature()
	if got != "unknown" {
		t.Errorf("expected unknown, got %q", got)
	}
}

func TestSignature_ArrayNoElem(t *testing.T) {
	p := ParamType{Kind: KindArray}
	got := p.Signature()
	if got != "any[]" {
		t.Errorf("expected any[], got %q", got)
	}
}

func TestSignature_FixedArrayNoElem(t *testing.T) {
	p := ParamType{Kind: KindFixedArr, Size: 3}
	got := p.Signature()
	if got != "any[]" {
		t.Errorf("expected any[], got %q", got)
	}
}

func TestAppendInt256Val_Negative_Cover(t *testing.T) {
	buf := make([]byte, 0, 32)
	v := big.NewInt(-42)
	result := appendInt256Val(buf, v)
	if len(result) != 32 {
		t.Fatal("expected 32 bytes")
	}
}

func TestAppendInt256Val_NegativeMax_Cover(t *testing.T) {
	buf := make([]byte, 0, 32)
	v := new(big.Int).Lsh(big.NewInt(-1), 255)
	result := appendInt256Val(buf, v)
	if len(result) != 32 {
		t.Fatal("expected 32 bytes")
	}
}

func TestAppendInt256Val_Nil_Cover(t *testing.T) {
	buf := make([]byte, 0, 32)
	result := appendInt256Val(buf, nil)
	if len(result) != 32 {
		t.Fatal("expected 32 bytes")
	}
}

func TestDynamicTailSize_NilTuple(t *testing.T) {
	tup := TTuple()
	s := dynamicTailSize(tup, nil)
	if s != 32 {
		t.Fatalf("expected 32 for nil, got %d", s)
	}
}

func TestDynamicTailSize_NilTupleEl(t *testing.T) {
	p := ParamType{Kind: KindTuple}
	s := dynamicTailSize(p, nil)
	if s != 32 {
		t.Fatalf("expected 32, got %d", s)
	}
}

func TestFixedArrayTailSize_StaticElems(t *testing.T) {
	p := TFixedArray(TUint256(), 3)
	vals := []any{NewU64(1), NewU64(2), NewU64(3)}
	s := fixedArrayTailSize(p, vals)
	if s != 96 {
		t.Fatalf("expected 96, got %d", s)
	}
}

func TestFixedArrayTailSize_DynamicElems(t *testing.T) {
	p := TFixedArray(TString(), 2)
	vals := []any{"hello", "world"}
	s := fixedArrayTailSize(p, vals)
	if s == 0 {
		t.Fatal("expected non-zero size")
	}
}

func TestFixedArrayTailSize_NoElem(t *testing.T) {
	p := ParamType{Kind: KindFixedArr, Size: 3}
	s := fixedArrayTailSize(p, nil)
	if s != 0 {
		t.Fatal("expected 0")
	}
}

func TestWadMul_WAD(t *testing.T) {
	a := WAD()
	b := WAD()
	result := WadMul(a, b)
	PutU256(result)
}

func TestWadDiv_Half(t *testing.T) {
	a := NewU64(1e18)
	b := NewU64(2)
	result := WadDiv(a, b)
	PutU256(result)
}
