package main

import (
	"log"
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"encoding/hex"
)

const DB_FILE = "blockchain.db"
const BLOCKS_BUCKET = "blocks"
const LAST_HASH = "last_hash"
const GENESIS_CORNBASE_DATA = "阿森纳淘汰了马竞"

// Blockchain 是一个 NewBlock 指针数组
type Blockchain struct {
	lastHash []byte
	db       *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		return false
	}

	return true
}

// 创建一个有创世块的新链
func NewBlockchain() *Blockchain {
	if dbExists() == false {
		fmt.Println("出错了, 区块链尚未初始化, 无法创建新区块")
		os.Exit(1)
	}
	var lastHash []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		lastHash = bucket.Get([]byte(LAST_HASH))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{lastHash, db}

	return &bc
}

// 创建一个新的区块链数据库
func InitBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("出错了, 不能重复初始化区块链")
		os.Exit(1)
	}

	var lastHash []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		coinbaseTransaction := CreateCoinbaseTransaction(address, GENESIS_CORNBASE_DATA)
		genisis := NewGenesisBlock(coinbaseTransaction)

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

		return nil
	})

	return &Blockchain{lastHash, db}
}

// MineBlock 向链中加入一个新块
// Data 在实际中就是交易
func (blockchain *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := blockchain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
		lastHash = bucket.Get([]byte(lastHash))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

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

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.lastHash, blockchain.db}
}

func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := iterator.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS_BUCKET))
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

func (blockchain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTransactions := blockchain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, transaction := range unspentTransactions {
		transactionID := hex.EncodeToString(transaction.ID)

		for index, output := range transaction.Outputs {
			if output.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += output.Value
				unspentOutputs[transactionID] = append(unspentOutputs[transactionID], index)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (blockchain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTransactions []Transaction
	mapSpentTransactionOutputIndex := make(map[string][]int)
	iterator := blockchain.Iterator()

	for {
		block := iterator.Next() //从最后一个区块, 从上往下遍历. 所以最后一块包含的输出肯定没花过??

		for _, transaction := range block.Transactions {
			transactionID := hex.EncodeToString(transaction.ID)

			//先看下当前交易中哪些输出可以被address花
		Outputs:
			for index, output := range transaction.Outputs {
				//首先检查当前输出是否已经花掉了
				if mapSpentTransactionOutputIndex[transactionID] != nil {
					for _, spentOutputIndex := range mapSpentTransactionOutputIndex[transactionID] {
						if spentOutputIndex == index { //当前outputIndex出现在已花费的outputIndex数组中, 所以当前输出已经被花掉
							continue Outputs
						}
					}
				}

				//如果当前输出没有被花掉, 而address又能解锁该输出, 那么就找到一笔address能花的交易
				if output.CanBeUnlockedWith(address) {
					unspentTransactions = append(unspentTransactions, *transaction)
				}
			}

			//再看下当前交易中哪些输入是被address花过的, 将相关信息(交易id, 输出的index)存起来. 等到遍历下一个区块时, 可以用来判断区块中的交易输出是否花了
			if transaction.IsCoinbaseTransaction() == false {
				for _, input := range transaction.Inputs {
					if input.CanUnlockOutputWith(address) { //谁有权解锁该输入, 就说明这笔钱已经被谁花掉了
						previousTransactionID := hex.EncodeToString(input.PreviousTransactionID)
						mapSpentTransactionOutputIndex[previousTransactionID] = append(mapSpentTransactionOutputIndex[previousTransactionID], input.OutputIndexInPreviousTransaction)
					}
				}
			}

		}

		if len(block.PreviousBlockHash) == 0 {
			break
		}
	}
	return unspentTransactions
}

func (blockchain *Blockchain) FindUTXO(address string) []TransactionOutput {
	var utxos []TransactionOutput
	unspentTransactions := blockchain.FindUnspentTransactions(address)

	for _, transaction := range unspentTransactions {
		for _, output := range transaction.Outputs {
			if output.CanBeUnlockedWith(address) {
				utxos = append(utxos, output)
			}
		}
	}

	return utxos
}