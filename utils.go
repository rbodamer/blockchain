package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

// PrintBlock prints a block
func PrintBlock(bc Block) {
	fmt.Printf("Prev Hash:  %x\n", bc.PrevBlockHash)
	fmt.Printf("Data:       %s\n", bc.Data)
	fmt.Printf("Hash:       %x\n", bc.Hash)
	pow := NewProofOfWork(&bc)
	fmt.Printf("PoW:        %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Println()
}
