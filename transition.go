package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction represents a transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return (len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].OutIndex == -1)
}

// SetHash sets the ID fields of the transaction
func (tx *Transaction) SetHash() {
	var result bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panicln(err)
	}

	hash = sha256.Sum256(result.Bytes())
	tx.ID = hash[:]
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, msg string) *Transaction {
	if msg == "" {
		msg = fmt.Sprintf("Reward to '%s'", to)
	}

	//TODO: why not set to nil
	txin := TXInput{[]byte{}, -1, nil, []byte(msg)}
	txout := NewTXOutput(subsidy, to)

	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetHash()

	return &tx
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amt int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets := GetWallets()
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PubKey)

	total, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amt)
	if total < amt {
		log.Panicln("not enough balance")
	}

	for id, outs := range validOutputs {
		txID, err := hex.DecodeString(id)
		if err != nil {
			log.Panicln(err)
		}

		for _, idx := range outs {
			input := TXInput{txID, idx, nil, wallet.PubKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTXOutput(amt, to))
	if total > amt {
		outputs = append(outputs, *NewTXOutput(total-amt, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetHash()

	return &tx
}
