package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// IntToHex represents an integer in the form of a byte slice
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panicln(err)
	}

	return buff.Bytes()
}
