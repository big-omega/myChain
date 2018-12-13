package main

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Blockchain is the transaction ledger
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// FindSpendableOutputs returns a collection of spendable transaction outputs for the given address and amount
func (bc *Blockchain) FindSpendableOutputs(address string, amt int) (int, map[string][]int) {
	spendableOutputs := make(map[string][]int)
	total := 0

	unspentTXs := bc.FindUnspentTransitions(address)

Collect:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for i, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				total += out.Value
				spendableOutputs[txID] = append(spendableOutputs[txID], i)

				if total >= amt {
					break Collect
				}
			}
		}
	}

	return total, spendableOutputs
}

// FindUnspentTransitions returns a list of unspent transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransitions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXO := make(map[string][]int)
	bci := bc.Iterator()
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for i, out := range tx.Vout {
				if spentTXO[txID] != nil {
					for _, idx := range spentTXO[txID] {
						if i == idx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}

				if !tx.IsCoinbase() {
					for _, in := range tx.Vin {
						if in.CanUnlockOutputWith(address) {
							inTxID := hex.EncodeToString(in.Txid)
							spentTXO[inTxID] = append(spentTXO[inTxID], in.OutIndex)
						}
					}
				}
			}
		}

		if len(bci.currentHash) == 0 {
			break
		}
	}
	return unspentTXs
}

// FindUTXOs finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXOs(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTXs := bc.FindUnspentTransitions(address)

	for _, tx := range unspentTXs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// Iterator returns a blockchain iterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(txs []*Transaction) {
	var tip []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panicln(err)
	}

	block := NewBlock(txs, tip)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err = b.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panicln(err)
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			log.Panicln(err)
		}

		return nil
	})

	bc.tip = block.Hash
}

// CreateBlockchain creates a new blockchain
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		log.Println("Blockchain already exists")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panicln(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		coinbase := NewCoinbaseTX(address, "")
		genesis := NewGenesisBlock(coinbase)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panicln(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panicln(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panicln(err)
		}

		tip = genesis.Hash

		return nil
	})
	if err != nil {
		log.Panicln(err)
	}

	return &Blockchain{tip, db}
}

// dbExists checks whether blockchain db exists
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// GetBlockchain returns the current blockchain instance
func GetBlockchain() *Blockchain {
	if !dbExists() {
		log.Panicln("No existing blockchain found, please Create one first")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panicln(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panicln(err)
	}

	return &Blockchain{tip, db}
}

// NewGenesisBlock generates a Genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
