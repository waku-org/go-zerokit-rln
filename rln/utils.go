package rln

import (
	"encoding/hex"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateWitness(
	identityCredential IdentityCredential, // TODO only the secret hash.
	data []byte,
	epoch [32]byte,
	merkleProof MerkleProof) RLNWitnessInput {

	return RLNWitnessInput{
		IdentityCredential: identityCredential,
		MerkleProof:        merkleProof,
		X:                  HashToBN255(data),
		Epoch:              epoch,
		RlnIdentifier:      RLN_IDENTIFIER,
	}
}

func ToIdentityCredentials(groupKeys [][]string) ([]IdentityCredential, error) {
	// groupKeys is  sequence of membership key tuples in the form of (identity key, identity commitment) all in the hexadecimal format
	// the toIdentityCredentials proc populates a sequence of IdentityCredentials using the supplied groupKeys
	// Returns an error if the conversion fails

	var groupIdCredentials []IdentityCredential

	for _, gk := range groupKeys {
		idTrapdoor, err := ToBytes32LE(gk[0])
		if err != nil {
			return nil, err
		}

		idNullifier, err := ToBytes32LE(gk[1])
		if err != nil {
			return nil, err
		}

		idSecretHash, err := ToBytes32LE(gk[2])
		if err != nil {
			return nil, err
		}

		idCommitment, err := ToBytes32LE(gk[3])
		if err != nil {
			return nil, err
		}

		groupIdCredentials = append(groupIdCredentials, IdentityCredential{
			IDTrapdoor:   idTrapdoor,
			IDNullifier:  idNullifier,
			IDSecretHash: idSecretHash,
			IDCommitment: idCommitment,
		})
	}

	return groupIdCredentials, nil
}

func Bytes32(b []byte) [32]byte {
	var result [32]byte
	copy(result[32-len(b):], b)
	return result
}

func Bytes128(b []byte) [128]byte {
	var result [128]byte
	copy(result[128-len(b):], b)
	return result
}

func Flatten(b [][32]byte) []byte {
	var result []byte
	for _, v := range b {
		result = append(result, v[:]...)
	}
	return result
}

func ToBytes32LE(hexStr string) ([32]byte, error) {

	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return [32]byte{}, err
	}

	bLen := len(b)
	for i := 0; i < bLen/2; i++ {
		b[i], b[bLen-i-1] = b[bLen-i-1], b[i]
	}

	return Bytes32(b), nil
}

func revert(b []byte) []byte {
	bLen := len(b)
	for i := 0; i < bLen/2; i++ {
		b[i], b[bLen-i-1] = b[bLen-i-1], b[i]
	}
	return b
}

// BigIntToBytes32 takes a *big.Int (which uses big endian) and converts it into a little endian 32 byte array
// Notice that is the *big.Int value contains an integer <= 2^248 - 1 (a 7 bytes value with all bits on), it will right-pad the result with 0s until
// the result has 32 bytes, i.e.:
// for a some bigInt whose `Bytes()` are {0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF}, using this function will return
// {0xEF, 0xCD, 0xAB, 0x90, 0x78, 0x56, 0x34, 0x12, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
func BigIntToBytes32(value *big.Int) [32]byte {
	b := revert(value.Bytes())
	tmp := make([]byte, 32)
	copy(tmp[0:len(b)], b)
	return Bytes32(tmp)
}

// Bytes32ToBigInt takes a little endian 32 byte array and returns a *big.Int (which uses big endian)
func Bytes32ToBigInt(value [32]byte) *big.Int {
	b := revert(value[:])
	result := new(big.Int)
	result.SetBytes(b)
	return result
}

// Hashes a byte array to a field element in BN254, as used by zerokit.
// Equivalent to: https://github.com/vacp2p/zerokit/blob/v0.3.4/rln/src/hashers.rs
func HashToBN255(data []byte) [32]byte {
	// Hash is fixed to 32 bytes
	hashed := crypto.Keccak256(data[:])

	// Convert to field element
	var frBN254 fr.Element
	frBN254.Unmarshal(revert(hashed))
	frBN254Bytes := frBN254.Bytes()

	// Return fixed size
	fixexLen := [32]byte{}
	copy(fixexLen[:], revert(frBN254Bytes[:]))
	return fixexLen
}
