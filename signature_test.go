package fastabi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeSignature(t *testing.T) {
	sig, err := DecodeSignature("transfer(address,uint256)")
	require.NoError(t, err)
	assert.NotEqual(t, [32]byte{}, sig)
}

func TestMustDecodeSignature(t *testing.T) {
	sig := MustDecodeSignature("transfer(address,uint256)")
	assert.NotEqual(t, [32]byte{}, sig)
}
