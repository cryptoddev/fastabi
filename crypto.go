package fastabi

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"unsafe"

	"github.com/holiman/uint256"
)

type Address [20]byte

func ParseAddress(s string) (Address, error) {
	var addr Address
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 40 {
		return addr, fmt.Errorf("invalid address length: %d", len(s))
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return addr, err
	}
	copy(addr[:], b)
	return addr, nil
}

func MustParseAddress(s string) Address {
	addr, err := ParseAddress(s)
	if err != nil {
		panic(err)
	}
	return addr
}

func (a Address) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a Address) String() string {
	return a.Hex()
}

func (a Address) IsZero() bool {
	return a == Address{}
}

func (a Address) Bytes() []byte {
	return a[:]
}

// U256 embeds uint256.Int inline (32 bytes, no pointers) so structs containing
// U256 value fields (not pointers) are GC-friendly (noscan).
type U256 struct {
	v uint256.Int
}

// u256Pool recycles heap-allocated U256 values to reduce GC pressure
// in hot paths (e.g., GetAmountOutV2 which calls GetU256 6+ times).
var u256Pool = sync.Pool{
	New: func() interface{} { return &U256{} },
}

func GetU256() *U256 {
	return u256Pool.Get().(*U256)
}

func PutU256(u *U256) {
	if u == nil {
		return
	}
	*u = U256{}
	u256Pool.Put(u)
}

func NewU64(v uint64) *U256 {
	u := &U256{}
	u.v.SetUint64(v)
	return u
}

func NewU256(v uint64) *U256 {
	return NewU64(v)
}

func MaxU256() *U256 {
	u := &U256{}
	u.v.SetAllOne()
	return u
}

func NewU256FromBig(b *big.Int) *U256 {
	u := &U256{}
	u.v.SetFromBig(b)
	return u
}

func NewU256FromHex(s string) (*U256, error) {
	u := &U256{}
	if err := u.v.SetFromHex(s); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *U256) SetUint64(v uint64) *U256 {
	u.v.SetUint64(v)
	return u
}

func (u *U256) Set(b *U256) *U256 {
	u.v.Set(&b.v)
	return u
}

func (u *U256) Inner() *uint256.Int {
	return &u.v
}

func (u *U256) ToBig() *big.Int {
	if u == nil {
		return big.NewInt(0)
	}
	return u.v.ToBig()
}

func (u *U256) Uint64() uint64 {
	if u == nil {
		return 0
	}
	return u.v.Uint64()
}

func (u *U256) IsZero() bool {
	if u == nil {
		return true
	}
	return u.v.IsZero()
}

func (u *U256) IsPositive() bool {
	if u == nil {
		return false
	}
	return !u.v.IsZero()
}

func (u *U256) Cmp(b *U256) int {
	return u.v.Cmp(&b.v)
}

func (u *U256) Eq(b *U256) bool {
	return u.v.Eq(&b.v)
}

func (u *U256) Lt(b *U256) bool {
	return u.v.Lt(&b.v)
}

func (u *U256) Gt(b *U256) bool {
	return u.v.Gt(&b.v)
}

func (u *U256) Add(b *U256) *U256 {
	u.v.Add(&u.v, &b.v)
	return u
}

func (u *U256) Sub(b *U256) *U256 {
	u.v.Sub(&u.v, &b.v)
	return u
}

func (u *U256) Mul(b *U256) *U256 {
	u.v.Mul(&u.v, &b.v)
	return u
}

func (u *U256) Div(b *U256) *U256 {
	u.v.Div(&u.v, &b.v)
	return u
}

func (u *U256) Rsh(n uint) *U256 {
	u.v.Rsh(&u.v, n)
	return u
}

// MulDiv calculates (u * b) / c with full precision, falling back to big.Int on overflow.
func (u *U256) MulDiv(b, c *U256) *U256 {
	result, overflow := new(uint256.Int).MulDivOverflow(&u.v, &b.v, &c.v)
	if overflow {
		ubig := u.v.ToBig()
		bbig := b.v.ToBig()
		cbig := c.v.ToBig()
		resultBig := new(big.Int).Mul(ubig, bbig)
		resultBig.Div(resultBig, cbig)
		u.v.SetFromBig(resultBig)
	} else {
		u.v = *result
	}
	return u
}

// MulDivRoundingUp calculates (u * b) / c with rounding up.
func (u *U256) MulDivRoundingUp(b, c *U256) *U256 {
	result, overflow := new(uint256.Int).MulDivOverflow(&u.v, &b.v, &c.v)
	if overflow {
		ubig := u.v.ToBig()
		bbig := b.v.ToBig()
		cbig := c.v.ToBig()
		product := new(big.Int).Mul(ubig, bbig)
		resultBig := new(big.Int).Div(product, cbig)
		remainder := new(big.Int).Mod(product, cbig)
		if remainder.Sign() > 0 {
			resultBig.Add(resultBig, big.NewInt(1))
		}
		u.v.SetFromBig(resultBig)
	} else {
		check := new(uint256.Int).Mul(result, &c.v)
		original := new(uint256.Int).Mul(&u.v, &b.v)
		u.v = *result
		if !check.Eq(original) {
			u.v.AddUint64(&u.v, 1)
		}
	}
	return u
}

func (u *U256) ToUint64() uint64 {
	if u == nil {
		return 0
	}
	return u.v.Uint64()
}

func (u *U256) SetBytes(b []byte) *U256 {
	u.v.SetBytes(b)
	return u
}

func (u *U256) Sqrt() *U256 {
	u.v.Sqrt(&u.v)
	return u
}

func (u *U256) String() string {
	if u == nil {
		return "0"
	}
	return u.v.ToBig().String()
}

func (u *U256) Hex() string {
	if u == nil {
		return "0x0"
	}
	return u.v.Hex()
}

func (u *U256) Bytes() []byte {
	if u == nil {
		return make([]byte, 32)
	}
	return u.v.Bytes()
}

func (u *U256) Bytes32() [32]byte {
	var b [32]byte
	if u == nil {
		return b
	}
	u.v.WriteToSlice(b[:])
	return b
}

func (u *U256) SetBytes32(b [32]byte) *U256 {
	u.v.SetBytes(b[:])
	return u
}

func (u *U256) Clone() *U256 {
	clone := &U256{}
	clone.v.Set(&u.v)
	return clone
}

func IntToU256(v int64) *U256 {
	u := &U256{}
	u.v.SetUint64(uint64(v))
	return u
}

// wadVal and halfWadVal are immutable — never modify these.
var (
	wadVal     uint256.Int
	halfWadVal uint256.Int
	rayVal     uint256.Int
)

func init() {
	wadVal.SetUint64(1e18)
	halfWadVal.SetUint64(5e17)
	rayVal.SetUint64(1)
	ten := new(uint256.Int).SetUint64(10)
	for i := 0; i < 27; i++ {
		rayVal.Mul(&rayVal, ten)
	}
}

func WAD() *U256 { return (*U256)(unsafe.Pointer(&wadVal)) }

func HALF_WAD() *U256 { return (*U256)(unsafe.Pointer(&halfWadVal)) }

func RAYInt() *uint256.Int { return &rayVal }

// WadMul performs wad multiplication: (a * b + 0.5*WAD) / WAD
func WadMul(a, b *U256) *U256 {
	result := GetU256()
	tmp := new(uint256.Int)
	tmp.Mul(&a.v, &b.v)
	tmp.Add(tmp, &halfWadVal)
	result.v.Div(tmp, &wadVal)
	return result
}

// WadDiv performs wad division: (a * WAD + b/2) / b
func WadDiv(a, b *U256) *U256 {
	result := GetU256()
	halfB := new(uint256.Int)
	halfB.Rsh(&b.v, 1)
	tmp := new(uint256.Int)
	tmp.Mul(&a.v, &wadVal)
	tmp.Add(tmp, halfB)
	result.v.Div(tmp, &b.v)
	return result
}

// WadLn computes natural logarithm for WAD fixed-point numbers.
func WadLn(x *U256) *U256 {
	if x.IsZero() || x.v.Lt(&wadVal) {
		return &U256{}
	}
	if x.v.Eq(&wadVal) {
		return &U256{}
	}

	u := new(uint256.Int).Sub(&x.v, &wadVal)

	u2 := new(uint256.Int).Mul(u, u)
	twoWAD := new(uint256.Int).Mul(&wadVal, new(uint256.Int).SetUint64(2))
	term2 := new(uint256.Int).Div(u2, twoWAD)

	u3 := new(uint256.Int).Mul(u2, u)
	threeWAD2 := new(uint256.Int).Mul(&wadVal, &wadVal)
	threeWAD2.Mul(threeWAD2, new(uint256.Int).SetUint64(3))
	term3 := new(uint256.Int).Div(u3, threeWAD2)

	result := &U256{}
	result.v.Set(u)
	if result.v.Cmp(term2) >= 0 {
		result.v.Sub(&result.v, term2)
	} else {
		result.v.SetUint64(0)
	}
	result.v.Add(&result.v, term3)

	return result
}

func ParseHexUint64(s string) (uint64, error) {
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}
	if s == "" {
		return 0, nil
	}

	var result uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			result = result*16 + uint64(c-'0')
		case c >= 'a' && c <= 'f':
			result = result*16 + uint64(c-'a'+10)
		case c >= 'A' && c <= 'F':
			result = result*16 + uint64(c-'A'+10)
		default:
			return 0, fmt.Errorf("invalid hex char: %c", c)
		}
	}
	return result, nil
}

func HexToBytes(s string) ([]byte, error) {
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}

	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		b, err := ParseHexByte(s[i], s[i+1])
		if err != nil {
			return nil, err
		}
		result[i/2] = b
	}
	return result, nil
}

func ParseHexByte(hi, lo byte) (byte, error) {
	var b byte
	switch {
	case hi >= '0' && hi <= '9':
		b = (hi - '0') << 4
	case hi >= 'a' && hi <= 'f':
		b = (hi - 'a' + 10) << 4
	case hi >= 'A' && hi <= 'F':
		b = (hi - 'A' + 10) << 4
	default:
		return 0, fmt.Errorf("invalid hex: %c", hi)
	}

	switch {
	case lo >= '0' && lo <= '9':
		b |= lo - '0'
	case lo >= 'a' && lo <= 'f':
		b |= lo - 'a' + 10
	case lo >= 'A' && lo <= 'F':
		b |= lo - 'A' + 10
	default:
		return 0, fmt.Errorf("invalid hex: %c", lo)
	}

	return b, nil
}

// AddressFromUint64 creates an Address from a uint64 value (big-endian, right-aligned).
// Useful in tests and benchmarks to generate unique addresses cheaply.
func AddressFromUint64(v uint64) Address {
	var a Address
	a[12] = byte(v >> 56)
	a[13] = byte(v >> 48)
	a[14] = byte(v >> 40)
	a[15] = byte(v >> 32)
	a[16] = byte(v >> 24)
	a[17] = byte(v >> 16)
	a[18] = byte(v >> 8)
	a[19] = byte(v)
	return a
}

func BytesToHex(b []byte) string {
	const hexDigits = "0123456789abcdef"
	result := make([]byte, len(b)*2)
	for i, v := range b {
		result[i*2] = hexDigits[v>>4]
		result[i*2+1] = hexDigits[v&0x0f]
	}
	return string(result)
}
