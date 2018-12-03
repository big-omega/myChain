package main

import (
	"fmt"
	"log"
	"strconv"
)

func init() {
	log.SetPrefix("DEBUG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

// main is the entry point for the program
func main() {
	fmt.Println("Welcome to my first blockchain")

	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to bigOmega")
	bc.AddBlock("Send 2 more BTC to bigOmega")

	for _, block := range bc.blocks {
		fmt.Printf("\nPrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PoW: %s", strconv.FormatBool(NewProofOfWork(block).Validate()))
		fmt.Println()
	}
}
