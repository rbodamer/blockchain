package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

const targetBits = 24

// Block keeps the block headers
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Blockchain is the actual blokchain
type Blockchain struct {
	blocks []*Block
}

// NewBlock creates and returns a new block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and return the first ever block in this chain, a Genesis block
func NewGenesisBlock() *Block {
	return NewBlock("Oasis", []byte{})
}

// PrintBlock prints a block
func PrintBlock(bc Block) {
	fmt.Printf("Prev Hash:  %x\n", bc.PrevBlockHash)
	fmt.Printf("Data:       %s\n", bc.Data)
	fmt.Printf("Hash:       %x\n", bc.Hash)
	pow := NewProofOfWork(&bc)
	fmt.Printf("Pow:        %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Println()
}

//NewBlockChain creates and returns a new blockhain
func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

//AddBlock adds a block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

//ProofOfWork contains block and target
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

//NewProofOfWork set the contract to determine POW
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}

	return pow
}

//prepareData merge block fields. Nonce is the counter.
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run Calculate proof of work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	maxNonce := math.MaxInt64
	nonce := 0

	fmt.Printf("Mining the block containing %s\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate - validate that a proof of work is correct
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

// IntToHex converts an int64 to a byte array
// func IntToHex(num int64) []byte {
// 	buff := new(bytes.Buffer)
// 	err := binary.Write(buff, binary.BigEndian, num)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return buff.Bytes()
// }

//Main - doit
func main() {
	bc := NewBlockChain()

	bc.AddBlock("Send 1BTC to Roger")
	bc.AddBlock("Send another BTC to Roger")

	for _, block := range bc.blocks {
		PrintBlock(*block)

	}
}
