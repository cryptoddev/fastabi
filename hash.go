package fastabi

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type Hash [32]byte

func ParseHash(s string) (Hash, error) {
	var h Hash
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 64 {
		return h, fmt.Errorf("invalid hash length: %d (want 64 hex chars)", len(s))
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return h, err
	}
	copy(h[:], b)
	return h, nil
}

func (h Hash) Hex() string {
	return "0x" + hex.EncodeToString(h[:])
}

func (h Hash) IsZero() bool {
	return h == Hash{}
}

func (h Hash) Bytes() []byte {
	return h[:]
}
