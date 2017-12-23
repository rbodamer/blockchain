package main

import (
	"bytes"
	"encoding/binary"
	"log"
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
