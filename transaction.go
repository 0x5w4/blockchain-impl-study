package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
)

type Transaction struct {
	Version  int32
	ID       []byte
	Vin      []TransactionInput
	Vout     []TransactionOutput
	LockTime uint32
	Fee      int64
}

func (tx *Transaction) Serialize() []byte {
	buf := make([]byte, 0, 256)

	tmp4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp4, uint32(tx.Version))
	buf = append(buf, tmp4...)

	buf = appendVarInt(buf, uint64(len(tx.Vin)))
	for _, vin := range tx.Vin {
		buf = append(buf, vin.PrevTxID...)

		binary.LittleEndian.PutUint32(tmp4, uint32(vin.PrevOutIndex))
		buf = append(buf, tmp4...)

		buf = appendVarBytes(buf, vin.Signature)
		buf = appendVarBytes(buf, vin.PubKey)
	}

	buf = appendVarInt(buf, uint64(len(tx.Vout)))
	for _, vout := range tx.Vout {
		tmp8 := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp8, uint64(vout.Value))
		buf = append(buf, tmp8...)

		buf = appendVarBytes(buf, vout.PubKeyHash)
	}

	binary.LittleEndian.PutUint32(tmp4, tx.LockTime)
	buf = append(buf, tmp4...)

	return buf
}

func (tx *Transaction) CalculateFee(prevTXs map[string]Transaction) int64 {
	var inputSum int64
	var outputSum int64

	for _, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.PrevTxID)]
		inputSum += int64(prevTx.Vout[vin.PrevOutIndex].Value)
	}

	for _, vout := range tx.Vout {
		outputSum += int64(vout.Value)
	}

	return inputSum - outputSum
}

func DeserializeTransaction(data []byte) Transaction {
	var tx Transaction
	r := bytes.NewReader(data)

	binary.Read(r, binary.LittleEndian, &tx.Version)

	vinCount, _ := decodeVarInt(r)
	for i := 0; i < int(vinCount); i++ {
		var vin TransactionInput

		vin.PrevTxID = make([]byte, 32)
		io.ReadFull(r, vin.PrevTxID)

		var vout uint32
		binary.Read(r, binary.LittleEndian, &vout)
		vin.PrevOutIndex = int(vout)

		sigLen, _ := decodeVarInt(r)
		vin.Signature = make([]byte, sigLen)
		io.ReadFull(r, vin.Signature)

		pubLen, _ := decodeVarInt(r)
		vin.PubKey = make([]byte, pubLen)
		io.ReadFull(r, vin.PubKey)

		tx.Vin = append(tx.Vin, vin)
	}

	voutCount, _ := decodeVarInt(r)
	for i := 0; i < int(voutCount); i++ {
		var vout TransactionOutput
		var val uint64
		binary.Read(r, binary.LittleEndian, &val)
		vout.Value = int(val)

		pubKeyHashLen, _ := decodeVarInt(r)
		vout.PubKeyHash = make([]byte, pubKeyHashLen)
		io.ReadFull(r, vout.PubKeyHash)

		tx.Vout = append(tx.Vout, vout)
	}

	binary.Read(r, binary.LittleEndian, &tx.LockTime)

	return tx
}

func DeserializeBlock(data []byte) *Block {
	r := bytes.NewReader(data)

	header, _ := DeserializeBlockHeaderFromReader(r)

	txCount, _ := decodeVarInt(r)

	var transactions []*Transaction
	for i := 0; i < int(txCount); i++ {

		tx := DeserializeTransactionFromReader(r)
		transactions = append(transactions, &tx)
	}

	return &Block{
		Header:       *header,
		Transactions: transactions,
		Height:       0,
	}
}

func DeserializeTransactionFromReader(r *bytes.Reader) Transaction {
	var tx Transaction

	binary.Read(r, binary.LittleEndian, &tx.Version)

	vinCount, _ := decodeVarInt(r)
	for i := 0; i < int(vinCount); i++ {
		var vin TransactionInput
		vin.PrevTxID = make([]byte, 32)
		io.ReadFull(r, vin.PrevTxID)

		var vout uint32
		binary.Read(r, binary.LittleEndian, &vout)
		vin.PrevOutIndex = int(vout)

		sigLen, _ := decodeVarInt(r)
		vin.Signature = make([]byte, sigLen)
		io.ReadFull(r, vin.Signature)

		pubLen, _ := decodeVarInt(r)
		vin.PubKey = make([]byte, pubLen)
		io.ReadFull(r, vin.PubKey)

		tx.Vin = append(tx.Vin, vin)
	}

	voutCount, _ := decodeVarInt(r)
	for i := 0; i < int(voutCount); i++ {
		var vout TransactionOutput
		var val uint64
		binary.Read(r, binary.LittleEndian, &val)
		vout.Value = int(val)

		pubKeyHashLen, _ := decodeVarInt(r)
		vout.PubKeyHash = make([]byte, pubKeyHashLen)
		io.ReadFull(r, vout.PubKeyHash)

		tx.Vout = append(tx.Vout, vout)
	}

	binary.Read(r, binary.LittleEndian, &tx.LockTime)

	return tx
}
