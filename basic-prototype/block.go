package main

import (
	"time"
	"fmt"
)

// Block 由区块头和交易两部分构成
// Timestamp, PrevBlockHash, Hash 属于区块头（block header）
// Timestamp     : 当前时间戳，也就是区块创建的时间
// PrevBlockHash : 前一个块的哈希
// Hash          : 当前块的哈希
// Data          : 区块实际存储的信息，比特币中也就是交易
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock 用于生成新块，参数需要 Data 与 PrevBlockHash
// 当前块的哈希会基于 Data 和 PrevBlockHash 计算得到
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Data:          []byte(data)}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	fmt.Printf("%c[1;0;32m%s%c[0m %s\n\n", 0x1B, "挖出来一块矿!", 0x1B, time.Now().Format("2006-01-02 15:04:05.000"))
	return block
}

// SetHash 设置当前块哈希
// Hash = sha256(PrevBlockHash + Data + Timestamp)
/*func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}*/

// NewGenesisBlock 生成创世块
func NewGenesisBlock() *Block {
	return NewBlock("创世块", []byte{})
}
