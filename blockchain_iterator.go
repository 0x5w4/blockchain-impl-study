package main

import (
	"log"

	"github.com/dgraph-io/badger/v4"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Database}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.CurrentHash)
		if err != nil {
			log.Panic(err)
		}

		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = append([]byte{}, val...)
			return nil
		})

		block = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.CurrentHash = block.Header.PrevBlockHash

	return block
}
