package rln

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigInt(t *testing.T) {
	base := big.NewInt(2)
	value := base.Exp(base, big.NewInt(248), nil)
	value = value.Sub(value, big.NewInt(1)) // 2^248 - 1

	b32Value := BigIntToBytes32(value)
	require.True(t, bytes.Equal(b32Value[:], []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0}))

	newValue := Bytes32ToBigInt(b32Value)
	require.True(t, bytes.Equal(newValue.Bytes(), value.Bytes()))
}

// TODO: Test Flatten
