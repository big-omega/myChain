package main

// Blockchain is the transaction ledger
type Blockchain struct {
	blocks []*Block
}

// AddBlock extends the blockchain
func (bc *Blockchain) AddBlock(data string) {
	newBlock := NewBlock(data, bc.blocks[len(bc.blocks)-1].Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewGenesisBlock generates a Genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain is a constructor for a blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
