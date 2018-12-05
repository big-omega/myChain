package main

import (
	"log"
)

// init set configurations for log information
func init() {
	log.SetPrefix("DEBUG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

// main is the entry point for the program
func main() {
	// fmt.Println("Welcome to my first blockchain")
	bc := NewBlockchain()
	defer bc.db.Close()

	cli := &CLI{bc}
	cli.Run()
}
