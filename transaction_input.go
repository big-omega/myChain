package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	OutIndex  int
	Signature []byte
	PubKey    []byte
}

// UseKey checks whether the input can be initiated by the public key
func (in *TXInput) UseKey(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash, HashPubKey(in.PubKey)) == 0
}
