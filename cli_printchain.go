package main

// printChain prints the blockchain
func (cli *CLI) printChain() {
	bc := GetBlockchain()
	defer bc.db.Close()

	bci := bc.Iterator()
	for {
		bci.Next().PrintBlock()

		if len(bci.currentHash) == 0 {
			break
		}
	}
}
