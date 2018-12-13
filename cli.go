package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// CLI is the interface to this program
type CLI struct{}

// printUsage prints out usage information
func (cli *CLI) printUsage() {
	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Println("Usage:")
	// fmt.Println("	addblock -data BLOCK_DATA: add a new block to the blockchain")
	fmt.Println("	createblockchain -address ADDRESS: create a blockchain with the given ADDRESS")
	fmt.Println("	getbalance -address ADDRESS: get balance of ADDRESS")
	fmt.Println("	printchain: print all the blocks in the blockchain")
	fmt.Println("	send: -from FROM -to TO -amount AMOUNT: send AMOUNT from FROM to TO")
	fmt.Println("-----------------------------------------------------------------------------------------------")
}

// getBalance returns the balance for the given account
func (cli *CLI) getBalance(address string) {
	bc := GetBlockchain()
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXOs(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("balance of %s: %d\n", address, balance)
}

// CreateBlockchain creates a blockchain with given address
func (cli *CLI) CreateBlockchain(address string) {
	bc := CreateBlockchain(address)
	defer bc.db.Close()

	fmt.Println("done creating blockchain")
}

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

// Run is the entry point for CLI
func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	// addBlockData := addBlockCmd.String("data", "", "block data")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "the address to send Genesis reward to")
	getBalanceAddress := getBalanceCmd.String("address", "", "the address to getbalance for")
	sendFrom := sendCmd.String("from", "", "source wallet address")
	sendTo := sendCmd.String("to", "", "destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "amount to send")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}

	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}

	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}

	default:
		cli.printUsage()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}

		cli.CreateBlockchain(*createBlockchainAddress)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

// send creates a transaction and the corresponding block to finish coin transferring
func (cli *CLI) send(from, to string, amt int) {
	bc := GetBlockchain()
	defer bc.db.Close()

	tx := NewTransaction(from, to, amt, bc)
	bc.MineBlock([]*Transaction{tx})
	log.Println("succeed sending")
}

// validateArgs checks whether the input is valid
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
