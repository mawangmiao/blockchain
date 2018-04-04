package main

func main() {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	commandLine := CommandLine{blockchain}
	commandLine.Run()
}
