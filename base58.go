package main

import (
	"bytes"
	"log"
	"math/big"
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte

	x := new(big.Int).SetBytes(input)
	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	// handle leading zero bytes
	for _, b := range input {
		if b == 0x00 {
			result = append(result, b58Alphabet[0])
		} else {
			break
		}
	}

	// reverse result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	base := big.NewInt(int64(len(b58Alphabet)))

	for _, b := range input {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		if charIndex < 0 {
			log.Panic("invalid Base58 character")
		}
		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()

	// restore leading zero bytes
	zeroCount := 0
	for _, b := range input {
		if b == b58Alphabet[0] {
			zeroCount++
		} else {
			break
		}
	}

	return append(bytes.Repeat([]byte{0x00}, zeroCount), decoded...)
}
