# fastabi — ABI Codec for EVM

A high-performance ABI encoder/decoder for EVM chains. No go-ethereum dependency.

## Features

- **ABIv2** — all Solidity types: `uint8`–`uint256`, `int8`–`int256`, `address`, `bool`, `bytes`, `bytes32`, `string`, arrays, nested tuples
- **Zero-alloc** — pooled `Decoder`/`Encoder` via `sync.Pool`, no heap allocations per call after warm-up
- **U256** — wraps `holiman/uint256` inline (no pointer, GC-friendly), with pooling for hot paths
- **Signature parser** — `ParseFunction("transfer(address,uint256)")` → name + parameter types
- **Keccak256** — 4-byte selectors from function signatures
- **Wad math** — `WadMul`, `WadDiv`, `WadLn` for DeFi fixed-point arithmetic

## Quick Start

```go
import "github.com/cryptoddev/fastabi"

// Encode a single uint256
data := fastabi.Encode1(fastabi.TUint256(), uint64(42))

// Decode
val, _ := fastabi.Decode1(fastabi.TUint256(), data)
u := val.(*fastabi.U256)

// Pooled low-level API for hot paths
enc := fastabi.GetEncoder()
defer fastabi.PutEncoder(enc)
enc.EncodeUint64(42)
data = enc.Bytes()
```

## Types

| Constructor | Solidity | Go Type |
|---|---|---|
| `TUint<N>()` | `uint8`–`uint256` | `uint64`, `*big.Int`, `*U256` |
| `TInt<N>()` | `int8`–`int256` | `int64`, `*big.Int` |
| `TAddress()` | `address` | `[20]byte`, `Address` |
| `TBool()` | `bool` | `bool` |
| `TBytes()` | `bytes` | `[]byte` |
| `TBytes32()` | `bytes32`/`bytesN` | `[32]byte` |
| `TString()` | `string` | `string` |
| `TArray(elem)` | `type[]` | `[]any` |
| `TFixedArray(elem, n)` | `type[n]` | `[]any` |
| `TTuple(...)` | `(a,b,c)` | `[]any` |

## Testing

```bash
go test ./...
go test -bench=. ./...
```

## License

MIT
