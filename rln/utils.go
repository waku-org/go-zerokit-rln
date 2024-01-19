package rln

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

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

func EndianConvertTODO(data []byte) [32]byte {

	hashGoEth := crypto.Keccak256(data[:])
	_ = hashGoEth
	//if len(hashGoEth) != 32 {
	//	fmt.Println("errorrrrrr")
	//	}
	myHash32 := [32]byte{}
	//	copy(myHash32[:], hashGoEth)

	// el hash esta controlado por ahora
	copy(myHash32[:], data)

	fmt.Println("inpit is: ", data)
	fmt.Println("hash is: ", myHash32)

	var uintVals [4]uint64

	for i := 0; i < 4; i++ {
		chunk := make([]byte, 8)
		copy(chunk, myHash32[i*8:(i+1)*8])
		fmt.Println("chunk is: ", chunk)

		myBig := new(big.Int)
		myBig.SetBytes(revert(chunk))
		fmt.Println("big is: ", myBig)
		uintVals[i] = myBig.Uint64()
	}

	fmt.Println("uintVals is: ", uintVals)

	returnthis := [32]byte{}

	return returnthis
}
