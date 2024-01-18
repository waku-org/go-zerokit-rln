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

	for _, testSize := range []int{0, 1, 8, 16, 20} {
		mProof := MerkleProof{
			PathElements: []MerkleNode{},
			PathIndexes:  []uint8{},
		}

		for i := 0; i < testSize; i++ {
			mProof.PathElements = append(mProof.PathElements, random32())
			mProof.PathIndexes = append(mProof.PathIndexes, uint8(i%2))
		}

		// Check the size is the expected
		ser := mProof.serialize()
		require.Equal(t, 8+testSize*32+testSize+8, len(ser))

		fmt.Println("path:", mProof.PathElements)
		fmt.Println("Serialized: ", ser)

		// Deserialize and check its matches the original
		desProof := MerkleProof{}
		err := desProof.deserialize(ser)
		require.NoError(t, err)
		require.Equal(t, mProof, desProof)
	}

	// TODO test for errors. eg different size.
}

func TestRLNWitnessInputSerDe(t *testing.T) {
	// TODO: Is data size arbitrary? or fixed to 32. How its decode in zerokit seems to be fixed to 32.
	// At least in the proof with custom witness function.
	data := [32]byte{0x00}
	depth := 20

	mProof := MerkleProof{
		PathElements: []MerkleNode{},
		PathIndexes:  []uint8{},
	}

	for i := 0; i < depth; i++ {
		mProof.PathElements = append(mProof.PathElements, random32())
		mProof.PathIndexes = append(mProof.PathIndexes, uint8(i%2))
	}

	witness := RLNWitnessInput{
		IdentityCredential: IdentityCredential{
			IDSecretHash: random32(),
		},
		MerkleProof:   mProof,
		Data:          data[:],
		Epoch:         ToEpoch(10),
		RlnIdentifier: [32]byte{0x00},
	}

	ser := witness.serialize()
	require.Equal(t, 32+8+depth*32+depth+8+32+32+32, len(ser))

	// TODO:
	//desWitness := RLNWitnessInput{}
	//err := desWitness.deserialize(ser)
}
