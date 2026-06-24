package fastabi

import (
	"golang.org/x/crypto/sha3"
)

func DecodeSignature(sig string) ([32]byte, error) {
	var topic0 [32]byte
	h := sha3.NewLegacyKeccak256()
	_, err := h.Write([]byte(sig))
	if err != nil {
		return topic0, err
	}
	sum := h.Sum(nil)
	copy(topic0[:], sum)
	return topic0, nil
}

func MustDecodeSignature(sig string) [32]byte {
	v, err := DecodeSignature(sig)
	if err != nil {
		panic(err)
	}
	return v
}
