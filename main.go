package main

import (
	"fmt"
)

// main is the entry point for the program
func main() {
	fmt.Println("Welcome to my first blockchain\n")

	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to bigOmega")
	bc.AddBlock("Send 2 more BTC to bigOmega")

	for _, block := range bc.blocks {
		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("")
	}
}
