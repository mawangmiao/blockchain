package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CommandLine struct {
}

const usage = `
Usage:
  print --- 打印出所有区块
  getbalance -address ADDRESS --- 查询地址余额
  initblockchain -address ADDRESS --- 初始化一个区块链并发放创世块奖励
  send -from FROM -to TO -amount AMOUNT --- 从FROM向TO转AMOUNT金额的钱
`

func (commandLine *CommandLine) initBlockchain(address string) {
	bc := InitBlockchain(address)
	bc.db.Close()
	fmt.Println("完成!")
}

func (commandLine *CommandLine) printUsage() {
	fmt.Println(usage)
}

func (commandLine *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		commandLine.printUsage()
		os.Exit(1)
	}
}

func (commandLine *CommandLine) Run() {
	commandLine.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	initblockchainCmd := flag.NewFlagSet("initblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "要查询余额的地址")
	initblockchainAddress := initblockchainCmd.String("address", "", "创世块奖励接受者的地址")
	sendFrom := sendCmd.String("from", "", "转账放的地址")
	sendTo := sendCmd.String("to", "", "接收方的地址")
	sendAmount := sendCmd.Int("amount", 0, "转账金额")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "initblockchain":
		err := initblockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		commandLine.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		commandLine.getBalance(*getBalanceAddress)
	}

	if printChainCmd.Parsed() {
		commandLine.PrintChain()
	}

	if initblockchainCmd.Parsed() {
		if *initblockchainAddress == "" {
			initblockchainCmd.Usage()
			os.Exit(1)
		}
		commandLine.initBlockchain(*initblockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		commandLine.send(*sendFrom, *sendTo, *sendAmount)
	}
}

func (commandLine *CommandLine) send(from, to string, amount int) {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	transaction := NewUTXOTransaction(from, to, amount, blockchain)
	blockchain.MineBlock([]*Transaction{transaction})
	fmt.Println("转账成功!")
}

func (commandLine *CommandLine) getBalance(address string) {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	balance := 0
	utxos := blockchain.FindUTXO(address)

	for _, output := range utxos {
		balance += output.Value
	}

	fmt.Printf("'%s'的余额为: %d\n", address, balance)
}

func (commandLine *CommandLine) PrintChain() {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	iterator := blockchain.Iterator()

	for {
		block := iterator.Next()

		fmt.Printf("前一块哈希: %x\n", block.PreviousBlockHash)
		fmt.Printf("哈希: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW验证结果: %s\n\n", strconv.FormatBool(pow.Validate()))

		if len(block.PreviousBlockHash) == 0 {
			break
		}
	}
}
