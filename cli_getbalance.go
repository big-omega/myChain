package main

import "fmt"

// getBalance returns the balance for the given account
func (cli *CLI) getBalance(address string) {
	bc := GetBlockchain()
	defer bc.db.Close()

	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXOs(pubKeyHash)

	balance := 0
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("balance of %s: %d\n", address, balance)
}
