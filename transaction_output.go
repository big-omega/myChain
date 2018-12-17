package main

import "bytes"

// TXOutput represents a transaction output
type TXOutput struct {
	Value int
	// if not capital, gob will ignore it when serializing
	PubKeyHash []byte
}

// IsLockedWithKey checks whether an output can be used by the owner of the public key
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput creates a new transaction output with the given value and address
func NewTXOutput(value int, address string) *TXOutput {
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	return &TXOutput{value, pubKeyHash}
}
