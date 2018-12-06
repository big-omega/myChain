package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// CLI is the interface to this program
type CLI struct {
	bc *Blockchain
}

// printUsage prints out usage information
func (cli *CLI) printUsage() {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Usage:")
	fmt.Println("	addblock -data BLOCK_DATA: add a new block to the blockchain")
	fmt.Println("	printchain: print all the blocks in the blockchain")
	fmt.Println("--------------------------------------------------------------------------------")
}

// validateArgs checks whether the input is valid
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// addBlock calls the function to add a new block
func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Printf("Add block succeeds\n")
}

// printChain calls the function to print chain information
func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()
		block.PrintBlock()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

}

// Run is the entry point for CLI
func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panicln(err)
		}

	default:
		cli.printUsage()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}

		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
