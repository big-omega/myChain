package main

import (
	"bytes"
	"math/big"
)

var base58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {
	var result []byte

	x := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)

	mod := new(big.Int)
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)
	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{base58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	base := big.NewInt(58)
	zeroBytes := 0

	for _, b := range input {
		if b == byte('1') {
			zeroBytes++
		} else {
			break
		}
	}

	for i := zeroBytes; i < len(input); i++ {
		idx := bytes.IndexByte(base58Alphabet, input[i])
		if idx == -1 {
			return []byte("")
		}
		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(idx)))
	}

	decoded := result.Bytes()
	return append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
}
