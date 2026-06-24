package fastabi

import (
	"testing"
)

func TestU256_BasicOperations(t *testing.T) {
	tests := []struct {
		name string
		a    uint64
		b    uint64
	}{
		{"small", 100, 200},
		{"medium", 1e12, 2e12},
		{"large", 1e18, 2e18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewU64(tt.a)
			b := NewU64(tt.b)

			sum := GetU256().Set(a)
			sum.Add(b)
			if sum.Uint64() != tt.a+tt.b {
				t.Errorf("Add: expected %d, got %d", tt.a+tt.b, sum.Uint64())
			}
			PutU256(sum)

			if tt.a > tt.b {
				diff := GetU256().SetUint64(tt.a)
				diff.Sub(NewU64(tt.b))
				if diff.Uint64() != tt.a-tt.b {
					t.Errorf("Sub: expected %d, got %d", tt.a-tt.b, diff.Uint64())
				}
				PutU256(diff)
			}

			product := GetU256().Set(a)
			product.Mul(b)
			expected := tt.a * tt.b
			if product.Uint64() != expected {
				t.Errorf("Mul: expected %d, got %d", expected, product.Uint64())
			}
			PutU256(product)

			if tt.b != 0 && tt.a <= 1e9 && tt.b <= 1e9 {
				quotient := GetU256().SetUint64(tt.a * tt.b)
				divisor := NewU64(tt.b)
				quotient.Div(divisor)
				if quotient.Uint64() != tt.a {
					t.Errorf("Div: expected %d, got %d", tt.a, quotient.Uint64())
				}
				PutU256(quotient)
				PutU256(divisor)
			}

			PutU256(a)
			PutU256(b)
		})
	}
}

func TestU256_Comparisons(t *testing.T) {
	tests := []struct {
		a  uint64
		b  uint64
		gt bool
		lt bool
		eq bool
	}{
		{100, 200, false, true, false},
		{200, 100, true, false, false},
		{100, 100, false, false, true},
	}

	for _, tt := range tests {
		a := NewU64(tt.a)
		b := NewU64(tt.b)

		if a.Gt(b) != tt.gt {
			t.Errorf("Gt(%d, %d): expected %v", tt.a, tt.b, tt.gt)
		}
		if a.Lt(b) != tt.lt {
			t.Errorf("Lt(%d, %d): expected %v", tt.a, tt.b, tt.lt)
		}
		if a.Eq(b) != tt.eq {
			t.Errorf("Eq(%d, %d): expected %v", tt.a, tt.b, tt.eq)
		}

		PutU256(a)
		PutU256(b)
	}
}

func TestU256_FixedPoint(t *testing.T) {
	wad := WAD()
	half := NewU64(5e17)

	result := GetU256().Set(wad)
	result.Mul(half)
	result.Div(WAD())

	expected := uint64(5e17)
	if result.Uint64() != expected {
		t.Errorf("Fixed point mul: expected %d, got %d", expected, result.Uint64())
	}

	PutU256(result)
	PutU256(half)
}

func TestU256_PoolReuse(t *testing.T) {
	u1 := GetU256()
	u1.SetUint64(12345)
	PutU256(u1)

	u2 := GetU256()
	u2.SetUint64(67890)
	PutU256(u2)
}

func TestU256_Clone(t *testing.T) {
	original := NewU64(12345)
	cloned := original.Clone()

	if !original.Eq(cloned) {
		t.Error("Clone should equal original")
	}

	cloned.Add(NewU64(1))

	if original.Eq(cloned) {
		t.Error("Clone modification should not affect original")
	}

	PutU256(original)
	PutU256(cloned)
}

func TestU256_MulDiv(t *testing.T) {
	x := NewU64(1000)
	y := NewU64(3000)
	z := NewU64(100)

	result := GetU256().Set(x)
	result.MulDiv(y, z)

	if result.Uint64() != 30000 {
		t.Errorf("MulDiv: expected 30000, got %d", result.Uint64())
	}

	PutU256(x)
	PutU256(y)
	PutU256(z)
	PutU256(result)
}



func TestU256_EdgeCases(t *testing.T) {
	t.Run("zero_add", func(t *testing.T) {
		zero := NewU64(0)
		one := NewU64(1)
		sum := GetU256().Set(zero)
		sum.Add(one)
		if sum.Uint64() != 1 {
			t.Errorf("0 + 1 = %d, want 1", sum.Uint64())
		}
		PutU256(sum)
		PutU256(zero)
		PutU256(one)
	})

	t.Run("zero_mul", func(t *testing.T) {
		zero := NewU64(0)
		five := NewU64(5)
		prod := GetU256().Set(zero)
		prod.Mul(five)
		if prod.Uint64() != 0 {
			t.Errorf("0 * 5 = %d, want 0", prod.Uint64())
		}
		PutU256(prod)
		PutU256(zero)
		PutU256(five)
	})

	t.Run("max_uint64", func(t *testing.T) {
		max := NewU64(^uint64(0))
		one := NewU64(1)
		sum := GetU256().Set(max)
		sum.Add(one)
		if sum.IsZero() {
			t.Error("Max + 1 should not be zero")
		}
		PutU256(sum)
		PutU256(max)
		PutU256(one)
	})

	t.Run("large_mul_div", func(t *testing.T) {
		a := GetU256()
		a.SetUint64(1e18)
		b := GetU256()
		b.SetUint64(1e18)
		c := GetU256()
		c.SetUint64(1e18)

		result := GetU256().Set(a)
		result.Mul(b)
		result.Div(c)

		if result.Uint64() != 1e18 {
			t.Errorf("Large mul/div: got %d, want %d", result.Uint64(), uint64(1e18))
		}
		PutU256(a)
		PutU256(b)
		PutU256(c)
		PutU256(result)
	})

	t.Run("nil_handling", func(t *testing.T) {
		var nilU256 *U256
		if !nilU256.IsZero() {
			t.Error("nil U256 should be IsZero")
		}
	})
}

func TestParseAddress(t *testing.T) {
	addr, err := ParseAddress("0x1234567890abcdef1234567890abcdef12345678")
	if err != nil {
		t.Fatalf("ParseAddress failed: %v", err)
	}
	if addr.Hex() != "0x1234567890abcdef1234567890abcdef12345678" {
		t.Errorf("Hex mismatch: %s", addr.Hex())
	}

	// Invalid length
	_, err = ParseAddress("0x1234")
	if err == nil {
		t.Error("expected error for short address")
	}

	// Invalid hex
	_, err = ParseAddress("0xGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG")
	if err == nil {
		t.Error("expected error for invalid hex")
	}

	// MustParseAddress panics on error
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustParseAddress should panic on invalid address")
			}
		}()
		MustParseAddress("0xinvalid")
	}()

	// MustParseAddress succeeds on valid address
	addr2 := MustParseAddress("0x0000000000000000000000000000000000000001")
	if addr2.IsZero() {
		t.Error("MustParseAddress should return non-zero address")
	}

	// Test Bytes() on address
	b := addr.Bytes()
	if len(b) != 20 {
		t.Errorf("Address.Bytes() length: got %d, want 20", len(b))
	}
}

func TestNewU256_Variants(t *testing.T) {
	// NewU256 alias
	a := NewU256(42)
	if a.Uint64() != 42 {
		t.Errorf("NewU256: got %d, want 42", a.Uint64())
	}
	PutU256(a)

	// MaxU256
	max := MaxU256()
	if max.IsZero() {
		t.Error("MaxU256 should not be zero")
	}
	PutU256(max)

	// NewU256FromHex
	h, err := NewU256FromHex("0xff")
	if err != nil || h.Uint64() != 255 {
		t.Errorf("NewU256FromHex: got %d, want 255 (err=%v)", h.Uint64(), err)
	}
	PutU256(h)

	// NewU256FromHex invalid
	_, err = NewU256FromHex("0xzz")
	if err == nil {
		t.Error("NewU256FromHex should error on invalid hex")
	}

	// IntToU256
	n := IntToU256(-42)
	// Uint64 of -42 as uint64 wraps
	n64 := n.Uint64()
	if n64 != ^uint64(41) {
		t.Errorf("IntToU256(-42): got %d, want %d", n64, ^uint64(41))
	}
	PutU256(n)
}

func TestU256_ConversionAndFormat(t *testing.T) {
	a := NewU64(1e18)

	// ToBig
	bigVal := a.ToBig()
	if bigVal.Uint64() != 1e18 {
		t.Errorf("ToBig: got %d, want %d", bigVal.Uint64(), uint64(1e18))
	}

	// String
	strVal := a.String()
	if strVal != "1000000000000000000" {
		t.Errorf("String: got %s, want 1000000000000000000", strVal)
	}

	// Hex
	hexVal := a.Hex()
	if hexVal != "0xde0b6b3a7640000" {
		t.Errorf("Hex: got %s, want 0xde0b6b3a7640000", hexVal)
	}

	// ToUint64
	n := a.ToUint64()
	if n != 1e18 {
		t.Errorf("ToUint64: got %d, want %d", n, uint64(1e18))
	}

	// Bytes32
	b32 := a.Bytes32()
	if b32[31] != 0x00 || b32[24] != 0x0d || b32[25] != 0xe0 {
		t.Errorf("Bytes32 unexpected: %x", b32)
	}

	PutU256(a)

	// Nil handling
	var nilU *U256
	if nilU.String() != "0" {
		t.Errorf("nil String: got %s, want 0", nilU.String())
	}
	if nilU.Hex() != "0x0" {
		t.Errorf("nil Hex: got %s, want 0x0", nilU.Hex())
	}
	if nilU.Bytes() == nil || len(nilU.Bytes()) != 32 {
		t.Errorf("nil Bytes: expected 32-byte slice")
	}
	b32nil := nilU.Bytes32()
	if b32nil != [32]byte{} {
		t.Errorf("nil Bytes32: expected zero array")
	}
	if nilU.ToUint64() != 0 {
		t.Errorf("nil ToUint64: expected 0")
	}
	if nilU.ToBig().Uint64() != 0 {
		t.Errorf("nil ToBig: expected 0")
	}
}

func TestU256_IsPositive(t *testing.T) {
	zero := NewU64(0)
	if zero.IsPositive() {
		t.Error("zero should not be IsPositive")
	}
	PutU256(zero)

	pos := NewU64(1)
	if !pos.IsPositive() {
		t.Error("1 should be IsPositive")
	}
	PutU256(pos)

	var nilU *U256
	if nilU.IsPositive() {
		t.Error("nil should not be IsPositive")
	}
}

func TestU256_Cmp(t *testing.T) {
	small := NewU64(100)
	large := NewU64(200)

	if small.Cmp(large) != -1 {
		t.Error("100 Cmp 200 should be -1")
	}
	if large.Cmp(small) != 1 {
		t.Error("200 Cmp 100 should be 1")
	}
	if small.Cmp(NewU64(100)) != 0 {
		t.Error("100 Cmp 100 should be 0")
	}

	PutU256(small)
	PutU256(large)
}

func TestU256_Rsh(t *testing.T) {
	a := NewU64(256) // 2^8
	a.Rsh(3)         // 256 >> 3 = 32
	if a.Uint64() != 32 {
		t.Errorf("256 >> 3 = %d, want 32", a.Uint64())
	}
	PutU256(a)
}

func TestU256_Sqrt(t *testing.T) {
	a := NewU64(100)
	a.Sqrt()
	if a.Uint64() != 10 {
		t.Errorf("sqrt(100) = %d, want 10", a.Uint64())
	}
	PutU256(a)

	// sqrt(0) = 0
	zero := NewU64(0)
	zero.Sqrt()
	if zero.Uint64() != 0 {
		t.Errorf("sqrt(0) = %d, want 0", zero.Uint64())
	}
	PutU256(zero)
}

func TestU256_MulDivRoundingUp(t *testing.T) {
	a := NewU64(100)
	b := NewU64(300)
	c := NewU64(100)
	a.MulDivRoundingUp(b, c) // (100 * 300) / 100 = 300
	if a.Uint64() != 300 {
		t.Errorf("MulDivRoundingUp: got %d, want 300", a.Uint64())
	}
	PutU256(a)
	PutU256(b)
	PutU256(c)

	// Rounding up: (2*3)/5 = 6/5 = 1.2 → rounds up to 2
	a2 := NewU64(2)
	b2 := NewU64(3)
	c2 := NewU64(5)
	a2.MulDivRoundingUp(b2, c2)
	if a2.Uint64() != 2 {
		t.Errorf("MulDivRoundingUp(2,3,5): got %d, want 2", a2.Uint64())
	}
	PutU256(a2)
	PutU256(b2)
	PutU256(c2)

	// Exact division: (10*5)/5 = 10
	a3 := NewU64(10)
	b3 := NewU64(5)
	c3 := NewU64(5)
	a3.MulDivRoundingUp(b3, c3)
	if a3.Uint64() != 10 {
		t.Errorf("MulDivRoundingUp(10,5,5): got %d, want 10", a3.Uint64())
	}
	PutU256(a3)
	PutU256(b3)
	PutU256(c3)
}

func TestWadMath(t *testing.T) {
	// WadMul: (1e18 * 2e18) / 1e18 = 2e18
	a := NewU64(1e18)
	b := NewU64(2e18)
	result := WadMul(a, b)
	if result.Uint64() != 2e18 {
		t.Errorf("WadMul: got %d, want %d", result.Uint64(), uint64(2e18))
	}
	PutU256(a)
	PutU256(b)
	PutU256(result)

	// WadDiv: (1e18 * 1e18) / 2e18 = 5e17
	a2 := NewU64(1e18)
	b2 := NewU64(2e18)
	result2 := WadDiv(a2, b2)
	if result2.Uint64() != 5e17 {
		t.Errorf("WadDiv: got %d, want %d", result2.Uint64(), uint64(5e17))
	}
	PutU256(a2)
	PutU256(b2)
	PutU256(result2)

	// WadLn(1e18) = 0 (ln(1) = 0)
	one := NewU64(1e18)
	lnOne := WadLn(one)
	if !lnOne.IsZero() {
		t.Errorf("WadLn(1e18): expected 0, got %d", lnOne.Uint64())
	}
	PutU256(one)
	PutU256(lnOne)

	// WadLn(0) = 0
	zero := NewU64(0)
	lnZero := WadLn(zero)
	if !lnZero.IsZero() {
		t.Errorf("WadLn(0): expected 0, got %d", lnZero.Uint64())
	}
	PutU256(zero)
	PutU256(lnZero)

	// WadLn(e) ≈ 1 (approximate check)
	eVal2 := NewU64(3 * 1e18) // ln(3) ≈ 1.0986 * 1e18
	ln3 := WadLn(eVal2)
	if ln3.IsZero() {
		t.Error("WadLn(3e18): expected non-zero")
	}
	PutU256(eVal2)
	PutU256(ln3)
}

func TestParseHexByte(t *testing.T) {
	b, err := ParseHexByte('a', 'f')
	if err != nil || b != 0xaf {
		t.Errorf("ParseHexByte('a','f'): got %x, want af", b)
	}

	_, err = ParseHexByte('z', 'f')
	if err == nil {
		t.Error("expected error for invalid high nibble")
	}

	_, err = ParseHexByte('a', 'z')
	if err == nil {
		t.Error("expected error for invalid low nibble")
	}
}

func TestBytesToHex(t *testing.T) {
	s := BytesToHex([]byte{0xde, 0xad, 0xbe, 0xef})
	if s != "deadbeef" {
		t.Errorf("BytesToHex: got %s, want deadbeef", s)
	}

	empty := BytesToHex([]byte{})
	if empty != "" {
		t.Errorf("BytesToHex empty: got %s, want ''", empty)
	}
}

func TestSetBytes32(t *testing.T) {
	var b [32]byte
	b[31] = 0x2a // value = 42

	u := GetU256()
	u.SetBytes32(b)
	if u.Uint64() != 42 {
		t.Errorf("SetBytes32: got %d, want 42", u.Uint64())
	}
	PutU256(u)
}

func TestParseHexUint64_AllCases(t *testing.T) {
	// No prefix
	n, err := ParseHexUint64("ff")
	if err != nil || n != 255 {
		t.Errorf("ParseHexUint64(ff): got %d, want 255", n)
	}

	// Empty
	n, err = ParseHexUint64("")
	if err != nil || n != 0 {
		t.Errorf("ParseHexUint64(''): got %d, want 0", n)
	}

	// Uppercase
	n, err = ParseHexUint64("0xFF")
	if err != nil || n != 255 {
		t.Errorf("ParseHexUint64(0xFF): got %d, want 255", n)
	}

	// Invalid char
	_, err = ParseHexUint64("0xGG")
	if err == nil {
		t.Error("expected error for 0xGG")
	}
}

func TestHexToBytes_AllCases(t *testing.T) {
	// Without prefix
	b, err := HexToBytes("deadbeef")
	if err != nil || len(b) != 4 {
		t.Errorf("HexToBytes(deadbeef): got %x, err=%v", b, err)
	}

	// Invalid hex
	_, err = HexToBytes("0xGG")
	if err == nil {
		t.Error("expected error for invalid hex")
	}

	// Empty with prefix
	b, err = HexToBytes("0x")
	if err != nil {
		t.Fatalf("HexToBytes(0x): %v", err)
	}
	if len(b) != 0 {
		t.Errorf("HexToBytes(0x): got len %d, want 0", len(b))
	}
}
