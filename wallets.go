package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Wallets stores a collection of wallet
type Wallets struct {
	Wallets map[string]*Wallet
}

const walletsFile = "wallet.dat"

// GetWallet returns a wallet by its address
func (ws *Wallets) GetWallet(address string) Wallet {
	wallet, exist := ws.Wallets[address]
	if !exist {
		log.Panicln("address not exist")
	}

	return *wallet
}

// LoadFromFile loads wallets from file
func (ws *Wallets) LoadFromFile() {
	walletContent, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		log.Panicln(err)
	}

	var wallets *Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(walletContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panicln(err)
	}

	ws.Wallets = wallets.Wallets
}

// NewWallet generates a new wallet and returns the cooresponding address
func (ws *Wallets) NewWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet
	return address
}

// SaveToFile stores the wallets data in a file
func (ws *Wallets) SaveToFile() {
	var result bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panicln(err)
	}

	err = ioutil.WriteFile(walletsFile, result.Bytes(), 0644)
	if err != nil {
		log.Panicln(err)
	}
}

// GetWallets gets the wallets from wallets file
func GetWallets() *Wallets {
	wallets := &Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	wallets.LoadFromFile()

	return wallets
}

// walletsExist checks whether the wallets file exists
func walletsExist() bool {
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		return false
	}

	return true
}
