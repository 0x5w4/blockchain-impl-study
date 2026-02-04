package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  addblock -data DATA - Add a block to the blockchain")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  createwallet - Create a new wallet")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet("node_1")
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}
}

func (cli *CLI) createBlockchain(address string) {

	chain := InitBlockchain(address, "node_1")
	chain.Close()
	fmt.Println("Done!")
}

func (cli *CLI) addBlock(data string) {

	chain := ContinueBlockchain("node_1")
	defer chain.Close()

	tx := NewCoinbaseTX("legacy_user", data)

	chain.AddBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	chain := ContinueBlockchain("node_1")
	defer chain.Close()

	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("============ Block %x ============\n", block.Header.Hash())
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.Header.PrevBlockHash)

		for _, tx := range block.Transactions {

			fmt.Printf("Tx: %x\n", tx.ID)
		}

		pow := NewProofOfWork(&block.Header)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.Header.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := NewWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)

	fmt.Printf("Your new address: %s\n", address)
}
