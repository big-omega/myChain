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
	fmt.Println("	createwallet: generate a new key pair and save it to wallets file")
	fmt.Println("	getbalance -address ADDRESS: get balance of ADDRESS")
	fmt.Println("	newaddress: generate new address")
	fmt.Println("	printchain: print all the blocks in the blockchain")
	fmt.Println("	send: -from FROM -to TO -amount AMOUNT: send AMOUNT from FROM to TO")
	fmt.Println("-----------------------------------------------------------------------------------------------")
}

// Run is the entry point for CLI
func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	newAddressCmd := flag.NewFlagSet("newaddress", flag.ExitOnError)
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
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}
	case "newaddress":
		err := newAddressCmd.Parse(os.Args[2:])
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

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceAddress)
	}

	if newAddressCmd.Parsed() {
		cli.newAddress()
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

// validateArgs checks whether the input is valid
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
