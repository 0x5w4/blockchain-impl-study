package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/badger/v4"
)

const (
	dbPath              = "./tmp/blocks_%s"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func DBExists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}

func InitBlockchain(address, nodeId string) *Blockchain {
	path := fmt.Sprintf("./tmp/blocks_%s", nodeId)
	if DBExists(path) {
		log.Panic("blockchain already exists")
	}

	var lastHash []byte

	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesisBlock := NewGenesisBlock(cbtx, 0x1d00ffff)

		fmt.Println("Genesis Block created")

		err = txn.Set(genesisBlock.Header.Hash(), genesisBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = txn.Set([]byte("l"), genesisBlock.Header.Hash())
		if err != nil {
			log.Panic(err)
		}

		lastHash = genesisBlock.Header.Hash()
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{lastHash, db}
}

func ContinueBlockchain(nodeId string) *Blockchain {
	path := fmt.Sprintf(dbPath, nodeId)
	if !DBExists(path) {
		log.Panic("blockchain not found")
	}

	var lastHash []byte

	opts := badger.DefaultOptions(path)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("l"))
		if err != nil {
			log.Panic(err)
		}
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{lastHash, db}
}

func (chain *Blockchain) AddBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(chain.LastHash)
		if err != nil {
			log.Panic(err)
		}

		var blockData []byte
		err = item.Value(func(val []byte) error {
			blockData = append([]byte{}, val...)
			return nil
		})

		lastBlock := DeserializeBlock(blockData)
		lastHash = lastBlock.Header.Hash()
		lastHeight = lastBlock.Height
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash, lastHeight+1, 0x1d00ffff)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Header.Hash(), newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = txn.Set([]byte("l"), newBlock.Header.Hash())
		if err != nil {
			log.Panic(err)
		}

		chain.LastHash = newBlock.Header.Hash()
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

func (chain *Blockchain) Close() {
	chain.Database.Close()
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxIn{
		PrevTxID:  []byte{},
		Vout:      0xffffffff,
		ScriptSig: []byte(data),
		Sequence:  0xffffffff,
	}

	txout := TxOut{
		Value:        10,
		ScriptPubKey: []byte(to),
	}

	tx := Transaction{Version: 1, Vin: []TxIn{txin}, Vout: []TxOut{txout}, LockTime: 0}
	return &tx
}
