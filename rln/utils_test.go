package rln

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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

func TestFlatten(t *testing.T) {
	in1 := [][32]byte{[32]byte{}}
	in2 := [][32]byte{[32]byte{0x00}, [32]byte{0x01}}
	in3 := [][32]byte{[32]byte{0x01, 0x02, 0x03}, [32]byte{0x04, 0x05, 0x06}}

	expected1 := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	expected2 := []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	expected3 := []byte{
		0x1, 0x2, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x4, 0x5, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	out1 := Flatten(in1)
	require.Equal(t, expected1, out1)

	out2 := Flatten(in2)
	require.Equal(t, expected2, out2)

	out3 := Flatten(in3)
	require.Equal(t, expected3, out3)
}

func TestTODO(t *testing.T) {
	// Inputs for proof generation
	msg := []byte{72, 7, 140, 254, 213, 99, 57, 234, 84, 150, 46, 114, 195, 124, 127, 88, 143, 196, 248, 229, 188, 23, 56, 39, 186, 117, 203, 16, 166, 58, 150, 165}

	conv := EndianConvertTODO(msg)

	fmt.Println(conv)

	ints := []uint64{16877630849297418056, 6376952776256034388, 2826034866254562447, 11931788747685459386}

	fmt.Println(ints[0])

	str := ""
	for i, _ := range ints {
		str = str + padBinaryString(uint64ToBinaryString(ints[4-i-1]), 64)
	}

	fmt.Println("-- ", str)

	byteArray, err := binaryStringToBytes(str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	myBigInt := new(big.Int)
	myBigInt.SetBytes(byteArray[:])

	fmt.Println("mybigint", myBigInt.String())

	fmt.Println(byteArray)

	// Expected
	// in := [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0,0, 0}
	// hashed: []byte{72, 7, 140, 254, 213, 99, 57, 234, 84, 150, 46, 114, 195, 124, 127, 88, 143, 196, 248, 229, 188, 23, 56, 39, 186, 117, 203, 16, 166, 58, 150, 165}
	// expected := [32]byte{69, 7, 140, 46, 26, 131, 147, 30, 161, 68, 2, 5, 234, 195, 227, 223, 119, 187, 116, 97, 153, 70, 71, 254, 60, 149, 54, 109, 77, 79, 105, 20}
}

func reverseString(input string) string {
	// Convert string to a slice of runes
	runes := []rune(input)

	// Reverse the order of runes
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Convert the slice of runes back to a string
	reversedString := string(runes)

	return reversedString
}

func binaryStringToBytes(binaryString string) ([32]byte, error) {
	var byteArray [32]byte

	// Ensure the binary string has 256 bits (32 bytes)
	if len(binaryString) != 256 {
		return byteArray, fmt.Errorf("binary string must have exactly 256 bits")
	}

	// Iterate over 32 chunks of 8 bits each and parse them to bytes
	for i := 0; i < 32; i++ {
		startIndex := i * 8
		endIndex := startIndex + 8
		bits := binaryString[startIndex:endIndex]

		// Parse the 8-bit chunk to a byte
		byteValue, err := strconv.ParseUint(bits, 2, 8)
		if err != nil {
			return byteArray, err
		}

		byteArray[i] = byte(byteValue)
	}

	return byteArray, nil
}

func padBinaryString(binaryString string, length int) string {
	// Calculate padding length
	paddingLength := length - len(binaryString)

	// Pad with zeros
	paddedBinaryString := strings.Repeat("0", paddingLength) + binaryString

	return paddedBinaryString
}

func bigIntToBinaryString(num *big.Int) string {
	return fmt.Sprintf("%b", num)
}

func uint64ToBinaryString(num uint64) string {
	return strconv.FormatUint(num, 2)
}
