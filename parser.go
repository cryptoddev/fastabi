package fastabi

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseType parses a human-readable ABI type string into a ParamType.
// Examples: "uint256", "address", "(uint8,address,bytes)[]", "string", "bytes32[3]"
func ParseType(s string) (ParamType, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return ParamType{}, fmt.Errorf("empty type")
	}
	return parseTypeStr(s)
}

func ParseTypes(s string) ([]ParamType, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	parts := splitTopLevel(s, ',')
	ts := make([]ParamType, len(parts))
	for i, p := range parts {
		t, err := parseTypeStr(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("param %d: %w", i, err)
		}
		ts[i] = t
	}
	return ts, nil
}

// ParseFunction parses a function signature like "transfer(address,uint256)"
// or "balanceOf(address):(uint256)" and returns the function name and input types.
func ParseFunction(sig string) (name string, inputs []ParamType, err error) {
	s := strings.TrimSpace(sig)
	paren := strings.Index(s, "(")
	if paren < 0 {
		return s, nil, nil
	}
	name = strings.TrimSpace(s[:paren])

	depth := 0
	end := -1
	for i := paren; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				end = i
				goto found
			}
		}
	}
found:
	if end < 0 {
		return name, nil, fmt.Errorf("unmatched parenthesis")
	}

	inputStr := s[paren+1 : end]
	inputs, err = ParseTypes(inputStr)
	return name, inputs, err
}

func parseTypeStr(s string) (ParamType, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return ParamType{}, fmt.Errorf("empty type")
	}

	if s[0] == '(' {
		return parseTupleType(s)
	}

	if idx := strings.LastIndex(s, "["); idx >= 0 && s[len(s)-1] == ']' {
		base := s[:idx]
		inner := s[idx+1 : len(s)-1]
		elem, err := parseTypeStr(base)
		if err != nil {
			return ParamType{}, err
		}
		if inner == "" {
			return TArray(elem), nil
		}
		n, err := strconv.Atoi(inner)
		if err != nil {
			return ParamType{}, fmt.Errorf("invalid array size: %s", inner)
		}
		return TFixedArray(elem, n), nil
	}

	switch s {
	case "uint", "uint256":
		return TUint256(), nil
	case "uint128":
		return TUint128(), nil
	case "uint64":
		return TUint64(), nil
	case "uint32":
		return TUint32(), nil
	case "uint24":
		return TUint24(), nil
	case "uint16":
		return TUint16(), nil
	case "uint8":
		return TUint8(), nil
	case "int", "int256":
		return TInt256(), nil
	case "int128":
		return TInt128(), nil
	case "int64":
		return TInt64(), nil
	case "int32":
		return TInt32(), nil
	case "int24":
		return TInt24(), nil
	case "int16":
		return TInt16(), nil
	case "int8":
		return TInt8(), nil
	case "address":
		return TAddress(), nil
	case "bool":
		return TBool(), nil
	case "bytes32":
		return TBytes32(), nil
	case "bytes":
		return TBytes(), nil
	case "string":
		return TString(), nil
	}

	if strings.HasPrefix(s, "bytes") && len(s) > 5 {
		n, err := strconv.Atoi(s[5:])
		if err == nil && n >= 1 && n <= 32 {
			return ParamType{Kind: KindBytes32, Size: n}, nil
		}
	}

	if strings.HasPrefix(s, "uint") && len(s) > 4 {
		n, err := strconv.Atoi(s[4:])
		if err == nil && n >= 8 && n <= 256 && n%8 == 0 {
			return ParamType{Kind: KindUint256, Size: n}, nil
		}
	}
	if strings.HasPrefix(s, "int") && len(s) > 3 {
		n, err := strconv.Atoi(s[3:])
		if err == nil && n >= 8 && n <= 256 && n%8 == 0 {
			return ParamType{Kind: KindInt256, Size: n}, nil
		}
	}

	return ParamType{}, fmt.Errorf("unknown type: %s", s)
}

func parseTupleType(s string) (ParamType, error) {
	if s[0] != '(' {
		return ParamType{}, fmt.Errorf("expected '(', got %c", s[0])
	}

	depth := 0
	end := -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				end = i
				goto foundTuple
			}
		}
	}
foundTuple:
	if end < 0 {
		return ParamType{}, fmt.Errorf("unmatched parenthesis in tuple")
	}

	inner := s[1:end]
	rest := s[end+1:]

	fields, err := ParseTypes(inner)
	if err != nil {
		return ParamType{}, fmt.Errorf("tuple fields: %w", err)
	}

	t := TTuple(fields...)

	if len(rest) > 0 && rest[0] == '[' {
		arrEnd := strings.Index(rest, "]")
		if arrEnd < 0 {
			return ParamType{}, fmt.Errorf("unmatched bracket after tuple")
		}
		innerArr := rest[1:arrEnd]
		if innerArr == "" {
			return TArray(t), nil
		}
		n, err := strconv.Atoi(innerArr)
		if err != nil {
			return ParamType{}, fmt.Errorf("invalid tuple array size: %s", innerArr)
		}
		return TFixedArray(t, n), nil
	}

	return t, nil
}

// splitTopLevel splits by separator, respecting parentheses and brackets.
func splitTopLevel(s string, sep rune) []string {
	var parts []string
	var buf strings.Builder
	depth := 0
	for _, r := range s {
		switch r {
		case '(', '[':
			depth++
			buf.WriteRune(r)
		case ')', ']':
			depth--
			buf.WriteRune(r)
		default:
			if r == sep && depth == 0 {
				parts = append(parts, buf.String())
				buf.Reset()
			} else {
				buf.WriteRune(r)
			}
		}
	}
	if buf.Len() > 0 || len(parts) == 0 {
		parts = append(parts, buf.String())
	}
	return parts
}

func ParamNames(ts []ParamType) []string {
	names := make([]string, len(ts))
	for i, t := range ts {
		if t.Name != "" {
			names[i] = t.Name
		} else {
			names[i] = t.Kind.String()
		}
	}
	return names
}

func (t ParamType) GoType() string {
	switch t.Kind {
	case KindUint256:
		return "*U256"
	case KindUint128, KindInt256, KindInt128:
		return "*big.Int"
	case KindUint64, KindUint32, KindUint24, KindUint16, KindUint8:
		return "uint64"
	case KindInt64, KindInt32, KindInt24, KindInt16, KindInt8:
		return "int64"
	case KindAddress:
		return "[20]byte"
	case KindBool:
		return "bool"
	case KindBytes32:
		return "[32]byte"
	case KindBytes:
		return "[]byte"
	case KindString:
		return "string"
	case KindArray:
		if t.Elem != nil {
			return "[]" + t.Elem.GoType()
		}
		return "[]any"
	case KindFixedArr:
		if t.Elem != nil {
			return fmt.Sprintf("[%d]%s", t.Size, t.Elem.GoType())
		}
		return "[]any"
	case KindTuple:
		return "[]any"
	default:
		return "any"
	}
}

func (t ParamType) Signature() string {
	switch t.Kind {
	case KindUint256:
		return "uint256"
	case KindUint128:
		return "uint128"
	case KindUint64:
		return "uint64"
	case KindUint32:
		return "uint32"
	case KindUint24:
		return "uint24"
	case KindUint16:
		return "uint16"
	case KindUint8:
		return "uint8"
	case KindInt256:
		return "int256"
	case KindInt128:
		return "int128"
	case KindInt64:
		return "int64"
	case KindInt32:
		return "int32"
	case KindInt24:
		return "int24"
	case KindInt16:
		return "int16"
	case KindInt8:
		return "int8"
	case KindAddress:
		return "address"
	case KindBool:
		return "bool"
	case KindBytes32:
		if t.Size > 0 {
			return fmt.Sprintf("bytes%d", t.Size)
		}
		return "bytes32"
	case KindBytes:
		return "bytes"
	case KindString:
		return "string"
	case KindArray:
		if t.Elem != nil {
			return t.Elem.Signature() + "[]"
		}
		return "any[]"
	case KindFixedArr:
		if t.Elem != nil {
			return fmt.Sprintf("%s[%d]", t.Elem.Signature(), t.Size)
		}
		return "any[]"
	case KindTuple:
		var parts []string
		for _, f := range t.TupleEl {
			parts = append(parts, f.Signature())
		}
		return "(" + strings.Join(parts, ",") + ")"
	default:
		return "unknown"
	}
}

func IsWhitespace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
