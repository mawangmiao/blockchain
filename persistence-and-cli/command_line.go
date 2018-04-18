package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CommandLine struct {
}

const usage = `
Usage:
  print                   打印出区块链的所有区块
  initblockchain -address ADDRESS - 初始化一个区块链并发放创世块奖励
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

	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	initblockchainCmd := flag.NewFlagSet("initblockchain", flag.ExitOnError)

	initblockchainAddress := initblockchainCmd.String("address", "", "创世块奖励接受者的地址")

	switch os.Args[1] {
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
}

func (commandLine *CommandLine) send(from, to string, amount int) {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	transaction := NewUTXOTransaction(from, to, amount, blockchain)
	blockchain.MineBlock([]*Transaction{transaction})
	fmt.Println("转账成功!")
}

func (commandLine *CommandLine) PrintChain() {
	/*	iterator := commandLine.blockchain.Iterator()

		for {
			block := iterator.Next()

			fmt.Printf("前一块哈希: %x\n", block.PreviousBlockHash)
			fmt.Printf("数据: %s\n", block.Transactions)
			fmt.Printf("哈希: %x\n", block.Hash)

			pow := NewProofOfWork(block)
			fmt.Printf("PoW验证结果: %s\n\n", strconv.FormatBool(pow.Validate()))

			if len(block.PreviousBlockHash) == 0 {
				break
			}
		}*/
}
