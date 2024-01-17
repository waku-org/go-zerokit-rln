package rln

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func random32() [32]byte {
	var randomBytes [32]byte
	_, _ = rand.Read(randomBytes[:])
	return randomBytes
}

func TestMerkleProofSerDe(t *testing.T) {

	mProof := MerkleProof{
		PathElements: []MerkleNode{},
		PathIndexes:  []uint8{},
	}

	ser := mProof.serialize()
	//require.Equal(t, []byte{0, 0, 0, 0}, ser, )
	require.Equal(t, 16, len(ser))

	mProof = MerkleProof{
		PathElements: []MerkleNode{[32]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0}},
		PathIndexes:  []uint8{0},
	}

	ser = mProof.serialize()
	//require.Equal(t, []byte{0, 0, 0, 0}, ser, )
	require.Equal(t, 49, len(ser))

	mProof = MerkleProof{}

	for i := 0; i < 16; i++ {
		mProof.PathElements = append(mProof.PathElements, random32())
		mProof.PathIndexes = append(mProof.PathIndexes, uint8(i%2))
	}

	ser = mProof.serialize()
	fmt.Println(ser)

	desProof := MerkleProof{}
	err := desProof.deserialize(ser)
	require.NoError(t, err)

	// TODO test for errors. eg different size.
}

func TestRLNWitnessInputSerDe(t *testing.T) {

}
