package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const checksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PriKey ecdsa.PrivateKey
	PubKey []byte
}

// Address returns the Base58 check address for the public key in the wallet
func (w *Wallet) Address() []byte {
	pubKeyHash := HashPubKey(w.PubKey)

	versionedPubKeyHash := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPubKeyHash)

	return Base58Encode(append(versionedPubKeyHash, checksum...))
}

// checksum generates a checksum for a versioned public key
func checksum(payload []byte) []byte {
	firstSha256 := sha256.Sum256(payload)
	secondSha256 := sha256.Sum256(firstSha256[:])
	return secondSha256[:checksumLen]
}

// HashPubKey hashes a public key
func HashPubKey(pubKey []byte) []byte {
	pubKeyHash256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(pubKeyHash256[:])
	if err != nil {
		log.Panicln(err)
	}

	return RIPEMD160Hasher.Sum(nil)
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	priKey, pubKey := newKeyPair()
	return &Wallet{priKey, pubKey}
}

// newKeyPair generates a new key pair
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicln(err)
	}
	pubKey := append(priKey.PublicKey.X.Bytes(), priKey.PublicKey.Y.Bytes()...)

	return *priKey, pubKey
}

// ValidateAddress checks whether a given address is valid
func ValidateAddress(address string) bool {
	decoded := Base58Decode([]byte(address))
	if len(decoded) == 0 {
		log.Println("invalid address")
		return false
	}

	version := decoded[:1]
	pubKeyHash := decoded[1 : len(decoded)-checksumLen]
	actualChecksum := decoded[len(decoded)-checksumLen:]

	targetChecksum := checksum(append(version, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}
