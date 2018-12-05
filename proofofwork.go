package main

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 15

// proofOfWork is an entity serving for Proof of Work
type proofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork serves as a constructor for a Proof of Work struct
func NewProofOfWork(b *Block) *proofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &proofOfWork{b, target}
}

// prepareData assembles the data to be hashed
func (pow *proofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run calculates the nonce for a block
func (pow *proofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		hash = sha256.Sum256(pow.prepareData(nonce))
		// fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}

	return nonce, hash[:]
}

// Validate checks the validity of a block concerning only the PoW
func (pow *proofOfWork) Validate() bool {
	var hashInt big.Int

	hash := sha256.Sum256(pow.prepareData(pow.block.Nonce))
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
