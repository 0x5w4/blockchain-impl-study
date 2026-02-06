package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
)

type TxIn struct {
	PrevTxID  []byte
	Vout      uint32
	ScriptSig []byte
	Sequence  uint32
}

type TxOut struct {
	Value        int64
	ScriptPubKey []byte
}

type Transaction struct {
	Version  int32
	Vin      []TxIn
	Vout     []TxOut
	LockTime uint32
}

func (tx *Transaction) Serialize() []byte {
	buf := make([]byte, 0, 256)

	tmp4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp4, uint32(tx.Version))
	buf = append(buf, tmp4...)

	buf = appendVarInt(buf, uint64(len(tx.Vin)))
	for _, vin := range tx.Vin {
		// PrevTxID should be 32 bytes on the wire
		if len(vin.PrevTxID) == 0 {
			buf = append(buf, make([]byte, 32)...)
		} else {
			buf = append(buf, vin.PrevTxID...)
		}

		binary.LittleEndian.PutUint32(tmp4, vin.Vout)
		buf = append(buf, tmp4...)

		buf = appendVarBytes(buf, vin.ScriptSig)

		binary.LittleEndian.PutUint32(tmp4, vin.Sequence)
		buf = append(buf, tmp4...)
	}

	buf = appendVarInt(buf, uint64(len(tx.Vout)))
	for _, vout := range tx.Vout {
		tmp8 := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp8, uint64(vout.Value))
		buf = append(buf, tmp8...)

		buf = appendVarBytes(buf, vout.ScriptPubKey)
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
		inputSum += prevTx.Vout[vin.Vout].Value
	}

	for _, vout := range tx.Vout {
		outputSum += vout.Value
	}

	return inputSum - outputSum
}

func DeserializeTransaction(data []byte) Transaction {
	r := bytes.NewReader(data)
	return DeserializeTransactionFromReader(r)
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
		var vin TxIn
		vin.PrevTxID = make([]byte, 32)
		io.ReadFull(r, vin.PrevTxID)

		binary.Read(r, binary.LittleEndian, &vin.Vout)

		scriptLen, _ := decodeVarInt(r)
		vin.ScriptSig = make([]byte, scriptLen)
		io.ReadFull(r, vin.ScriptSig)

		binary.Read(r, binary.LittleEndian, &vin.Sequence)

		tx.Vin = append(tx.Vin, vin)
	}

	voutCount, _ := decodeVarInt(r)
	for i := 0; i < int(voutCount); i++ {
		var vout TxOut
		binary.Read(r, binary.LittleEndian, &vout.Value)

		scriptLen, _ := decodeVarInt(r)
		vout.ScriptPubKey = make([]byte, scriptLen)
		io.ReadFull(r, vout.ScriptPubKey)

		tx.Vout = append(tx.Vout, vout)
	}

	binary.Read(r, binary.LittleEndian, &tx.LockTime)

	return tx
}
