package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Block is a container for data in blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Nonce         int
	Hash          []byte
}

// HashTransactions returns the hash of all transactions in the block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var hash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	hash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return hash[:]
}

// NewBlock is a constructor for a block
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), txs, prevBlockHash, 0, []byte{}}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panicln(err)
	}

	return result.Bytes()
}

// PrintBlock prints the block info
func (b *Block) PrintBlock() {
	fmt.Printf("PrevHash: %x\n", b.PrevBlockHash)
	fmt.Printf("Hash: %x\n", b.Hash)
	for _, tx := range b.Transactions {
		fmt.Printf("Input: %v\n", tx.Vin)
		fmt.Printf("Output: %v\n", tx.Vout)
	}
	fmt.Printf("PoW: %s\n", strconv.FormatBool(NewProofOfWork(b).Validate()))
	fmt.Println()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panicln(err)
	}

	return &block
}
