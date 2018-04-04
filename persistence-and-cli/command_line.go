package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CommandLine struct {
	blockchain *Blockchain
}

const usage = `
Usage:
  add -Data BLOCK_DATA    add a block to the blockchain
  print                   print all the blocks of the blockchain
`

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

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
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

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		commandLine.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		commandLine.PrintChain()
	}
}

func (commandLine *CommandLine) AddBlock(data string) {
	commandLine.blockchain.AddBlock(data)
}

func (commandLine *CommandLine) PrintChain() {
	iterator := commandLine.blockchain.Iterator()

	for {
		block := iterator.Next()

		fmt.Printf("前一块哈希: %x\n", block.PreviousBlockHash)
		fmt.Printf("数据: %s\n", block.Data)
		fmt.Printf("哈希: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW验证结果: %s\n\n", strconv.FormatBool(pow.Validate()))

		if len(block.PreviousBlockHash) == 0 {
			break
		}
	}
}
