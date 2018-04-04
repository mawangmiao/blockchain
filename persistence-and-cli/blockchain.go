package main

import (
	"log"
	"github.com/boltdb/bolt"
)

const db_file = "blockchain.db"
const blocks_bucket = "blocks"
const last_hash = "last_hash"

// Blockchain 是一个 Block 指针数组
type Blockchain struct {
	lastHash []byte
	db       *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// NewBlockchain 创建一个有创世块的链
func NewBlockchain() *Blockchain {
	var lastHash []byte
	db, err := bolt.Open(db_file, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocks_bucket))

		if bucket == nil {
			genisis := NewGenesisBlock()
			bucket, err := tx.CreateBucket([]byte(blocks_bucket))
			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put(genisis.Hash, genisis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put([]byte(last_hash), genisis.Hash)

			if err != nil {
				log.Panic(err)
			}
			lastHash = genisis.Hash
		} else {
			lastHash = bucket.Get([]byte(last_hash))
		}

		return nil
	})

	return &Blockchain{lastHash, db}
}

// AddBlock 向链中加入一个新块
// Data 在实际中就是交易
func (blockchain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := blockchain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocks_bucket))
		lastHash = bucket.Get([]byte(last_hash))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocks_bucket))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte(last_hash), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		blockchain.lastHash = newBlock.Hash
		return nil
	})
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.lastHash, blockchain.db}
}

func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := iterator.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocks_bucket))
		encodedBlock := bucket.Get(iterator.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	iterator.currentHash = block.PreviousBlockHash
	return block
}
