package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// BlockchainIterator is used for iteratoring the blockchain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns the next block starting from the tip
func (bci *BlockchainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		block = DeserializeBlock(b.Get(bci.currentHash))

		return nil
	})
	if err != nil {
		log.Panicln(err)
	}

	bci.currentHash = block.PrevBlockHash

	return block
}
