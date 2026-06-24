package fastabi

type Kind uint8

const (
	KindUnknown Kind = iota
	KindUint8
	KindUint16
	KindUint24
	KindUint32
	KindUint64
	KindUint128
	KindUint256
	KindInt8
	KindInt16
	KindInt24
	KindInt32
	KindInt64
	KindInt128
	KindInt256
	KindAddress
	KindBool
	KindBytes32
	KindBytes
	KindString
	KindArray
	KindFixedArr
	KindTuple
)

func (k Kind) String() string {
	switch k {
	case KindUint8:
		return "uint8"
	case KindUint16:
		return "uint16"
	case KindUint24:
		return "uint24"
	case KindUint32:
		return "uint32"
	case KindUint64:
		return "uint64"
	case KindUint128:
		return "uint128"
	case KindUint256:
		return "uint256"
	case KindInt8:
		return "int8"
	case KindInt16:
		return "int16"
	case KindInt24:
		return "int24"
	case KindInt32:
		return "int32"
	case KindInt64:
		return "int64"
	case KindInt128:
		return "int128"
	case KindInt256:
		return "int256"
	case KindAddress:
		return "address"
	case KindBool:
		return "bool"
	case KindBytes32:
		return "bytes32"
	case KindBytes:
		return "bytes"
	case KindString:
		return "string"
	case KindArray:
		return "array"
	case KindFixedArr:
		return "fixedArray"
	case KindTuple:
		return "tuple"
	default:
		return "unknown"
	}
}

type ParamType struct {
	Kind    Kind
	Size    int        // bit-width for uint/int, length for fixedArray/fixedBytes
	Elem    *ParamType
	TupleEl []ParamType
	Name    string
}

func (t ParamType) IsDynamic() bool {
	switch t.Kind {
	case KindBytes, KindString, KindArray:
		return true
	case KindTuple:
		for _, e := range t.TupleEl {
			if e.IsDynamic() {
				return true
			}
		}
		return false
	case KindFixedArr:
		if t.Elem == nil {
			return false
		}
		return t.Elem.IsDynamic()
	default:
		return false
	}
}

func (t ParamType) StaticSize() int {
	if t.IsDynamic() {
		return 0
	}
	switch t.Kind {
	case KindTuple:
		size := 0
		for _, e := range t.TupleEl {
			size += e.StaticSize()
		}
		return size
	case KindFixedArr:
		if t.Elem == nil {
			return 0
		}
		return t.Size * t.Elem.StaticSize()
	default:
		return 32
	}
}

func TUint256() ParamType { return ParamType{Kind: KindUint256} }
func TUint128() ParamType { return ParamType{Kind: KindUint128} }
func TUint64() ParamType  { return ParamType{Kind: KindUint64} }
func TUint32() ParamType  { return ParamType{Kind: KindUint32} }
func TUint24() ParamType  { return ParamType{Kind: KindUint24} }
func TUint16() ParamType  { return ParamType{Kind: KindUint16} }
func TUint8() ParamType   { return ParamType{Kind: KindUint8} }
func TInt256() ParamType  { return ParamType{Kind: KindInt256} }
func TInt128() ParamType  { return ParamType{Kind: KindInt128} }
func TInt64() ParamType   { return ParamType{Kind: KindInt64} }
func TInt32() ParamType   { return ParamType{Kind: KindInt32} }
func TInt24() ParamType   { return ParamType{Kind: KindInt24} }
func TInt16() ParamType   { return ParamType{Kind: KindInt16} }
func TInt8() ParamType    { return ParamType{Kind: KindInt8} }
func TAddress() ParamType { return ParamType{Kind: KindAddress} }
func TBool() ParamType    { return ParamType{Kind: KindBool} }
func TBytes32() ParamType { return ParamType{Kind: KindBytes32} }
func TBytes() ParamType   { return ParamType{Kind: KindBytes} }
func TString() ParamType  { return ParamType{Kind: KindString} }

func TArray(elem ParamType) ParamType {
	return ParamType{Kind: KindArray, Elem: &elem}
}

func TFixedArray(elem ParamType, n int) ParamType {
	return ParamType{Kind: KindFixedArr, Elem: &elem, Size: n}
}

func TTuple(fields ...ParamType) ParamType {
	return ParamType{Kind: KindTuple, TupleEl: fields}
}

func Named(name string, t ParamType) ParamType {
	t.Name = name
	return t
}
