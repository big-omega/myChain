package main

import (
	"fmt"
	"log"
)

// CreateBlockchain creates a blockchain with given address
func (cli *CLI) CreateBlockchain(address string) {
	if !ValidateAddress(address) {
		log.Panicln("invalid address")
	}

	bc := CreateBlockchain(address)
	defer bc.db.Close()

	fmt.Println("create blockchain success!")
}
