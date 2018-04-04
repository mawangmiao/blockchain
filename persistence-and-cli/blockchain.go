package main

import (
	"log"
	"github.com/boltdb/bolt"
)

const DB_FILE = "blockchain.db"
const BLOCKS_BUCKET = "blocks"
const LAST_HASH = "last_hash"

// Blockchain 是一个 Block 指针数组
type Blockchain struct {
	lastHash []byte
	db       *bolt.DB
}

// NewBlockChain 创建一个有创世块的链
func NewBlockChain() *Blockchain {
	var lastHash []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))

		if bucket == nil {
			genisis := NewGenesisBlock()
			bucket, err := tx.CreateBucket([]byte(BLOCKS_BUCKET))
			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put(genisis.Hash, genisis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = bucket.Put([]byte(LAST_HASH), genisis.Hash)

			if err != nil {
				log.Panic(err)
			}
			lastHash = genisis.Hash
		} else {
			lastHash = bucket.Get([]byte(LAST_HASH))
		}

		return nil
	})

	return &Blockchain{lastHash, db}
}

// AddBlock 向链中加入一个新块
// data 在实际中就是交易
func (blockchain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := blockchain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		lastHash = bucket.Get([]byte(LAST_HASH))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte(LAST_HASH), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		blockchain.lastHash = newBlock.Hash
		return nil
	})
}
