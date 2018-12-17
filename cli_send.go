package main

import "log"

// send creates a transaction and the corresponding block to finish coin transferring
func (cli *CLI) send(from, to string, amt int) {
	bc := GetBlockchain()
	defer bc.db.Close()

	tx := NewTransaction(from, to, amt, bc)
	bc.MineBlock([]*Transaction{tx})
	log.Println("send success")
}
