package main

import (
	"fmt"
	"log"
	"os"
)

// newAddress generates a new address
func (cli *CLI) newAddress() {
	if !walletsExist() {
		log.Println("wallet not exist, please create one first")
		os.Exit(1)
	}

	ws := Wallets{}
	ws.Wallets = make(map[string]*Wallet)
	ws.LoadFromFile()

	address := ws.NewWallet()
	ws.SaveToFile()
	fmt.Printf("new address: %s\n", address)
}
