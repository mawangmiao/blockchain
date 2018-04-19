package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

const COINBASE_DEFAULT_OUTPUT_INDEX = -1
const SUBSIDY = 50

type TransactionInput struct {
	PreviousTransactionID            []byte
	OutputIndexInPreviousTransaction int
	ScriptSig                        string
}

type TransactionOutput struct {
	Value        int
	ScriptPubKey string
}

type Transaction struct {
	ID      []byte
	Inputs  []TransactionInput
	Outputs []TransactionOutput
}

func (transaction Transaction) IsCoinbaseTransaction() bool {
	return len(transaction.Inputs) == 1 && len(transaction.Inputs[0].PreviousTransactionID) == 0 && transaction.Inputs[0].OutputIndexInPreviousTransaction == COINBASE_DEFAULT_OUTPUT_INDEX
}

func CreateCoinbaseTransaction(to, data string) *Transaction {
	if data == "" {
		data = "奖励给 " + to
	}
	transactionInput := TransactionInput{[]byte{}, COINBASE_DEFAULT_OUTPUT_INDEX, data}
	transactionOutput := TransactionOutput{SUBSIDY, to}
	transaction := Transaction{nil, []TransactionInput{transactionInput}, []TransactionOutput{transactionOutput}}
	transaction.SetID()
	return &transaction
}

func (transaction *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(transaction)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	transaction.ID = hash[:]
}

func (transactionInput *TransactionInput) CanUnlockOutputWith(unlockingData string) bool {
	return transactionInput.ScriptSig == unlockingData
}

func (transactionOutput *TransactionOutput) CanBeUnlockedWith(unlockingData string) bool {
	return transactionOutput.ScriptPubKey == unlockingData
}

func NewUTXOTransaction(from, to string, amount int, blockchain *Blockchain) *Transaction {
	var inputs []TransactionInput
	var outputs []TransactionOutput

	outputValueSum, spendableOutputs := blockchain.FindSpendableOutputs(from, amount)
	if outputValueSum < amount {
		log.Panic("余额不足")
	}

	//创建交易的输入
	for encodedTransactionID, outputIndexes := range spendableOutputs {
		transactionID, err := hex.DecodeString(encodedTransactionID)
		if (err != nil) {
			log.Panic(err)
		}
		for _, outputIndex := range outputIndexes {
			input := TransactionInput{transactionID, outputIndex, from}
			inputs = append(inputs, input)
		}
	}
	//创建交易的输出
	outputs = append(outputs, TransactionOutput{amount, to})
	if outputValueSum > amount { //找零
		outputs = append(outputs, TransactionOutput{outputValueSum - amount, from})
	}

	newTransaction := Transaction{nil, inputs, outputs,}
	newTransaction.SetID() //TODO 内部逻辑不应在外面显示调用, 应移到初始化函数中去
	return &newTransaction
}
