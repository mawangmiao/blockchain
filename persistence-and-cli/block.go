package main

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

// NewBlock 由区块头和交易两部分构成
// Timestamp, PreviousBlockHash, Hash 属于区块头（block header）
// Timestamp     : 当前时间戳，也就是区块创建的时间
// PreviousBlockHash : 前一个块的哈希
// Hash          : 当前块的哈希
// Data          : 区块实际存储的信息，比特币中也就是交易
type Block struct {
	Timestamp         int64
	Transactions      []*Transaction
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
}

// 用于生成新块，参数需要 Data 与 PreviousBlockHash
// 当前块的哈希会基于 Data 和 PreviousBlockHash 计算得到
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	fmt.Printf("%c[1;0;32m%s%c[0m %s\n\n", 0x1B, "挖出来一块矿!", 0x1B, time.Now().Format("2006-01-02 15:04:05.000"))
	return block
}

// 生成创世块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
