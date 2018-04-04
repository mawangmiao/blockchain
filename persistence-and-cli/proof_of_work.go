package main

import (
	"math/big"
	"bytes"
	"fmt"
	"math"
	"crypto/sha256"
	"time"
)

const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1);
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{pow.block.PreviousBlockHash, pow.block.Data, IntToHex(pow.block.Timestamp), IntToHex(int64(targetBits)), IntToHex(int64(nonce))},
		[]byte{})
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("开始挖矿, 数据为 \"%s\" %s\n", pow.block.Data, time.Now().Format("2006-01-02 15:04:05.000"))
	for nonce < math.MaxInt64 {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.PrepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
