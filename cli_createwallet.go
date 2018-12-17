// The wallet here represents an unified interface to user, although it represents a Wallets internally.
// By default, the create action will generate a new key pair (as well as the corresponding address)

package main

import (
	"fmt"
	"log"
	"os"
)

// createWallet creates a new wallets file, and prints the first address generated
func (cli *CLI) createWallet() {
	if walletsExist() {
		log.Println("wallet already exists")
		os.Exit(1)
	}

	if _, err := os.Create(walletsFile); err != nil {
		log.Panicln(err)
	}

	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	address := wallets.NewWallet()
	wallets.SaveToFile()
	fmt.Printf("create wallet success. your address is: %s\n", address)
}
