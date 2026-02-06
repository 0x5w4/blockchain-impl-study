package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type BlockHeader struct {
	Version       uint32
	PrevBlockHash []byte
	MerkleRoot    []byte
	Timestamp     uint32
	Nonce         uint32
	Bits          uint32
}

func (h *BlockHeader) Serialize() []byte {
	buf := make([]byte, 0, 80) // Bitcoin header = 80 bytes

	tmp := make([]byte, 4)

	binary.LittleEndian.PutUint32(tmp, h.Version)
	buf = append(buf, tmp...)

	buf = append(buf, h.PrevBlockHash...)
	buf = append(buf, h.MerkleRoot...)

	binary.LittleEndian.PutUint32(tmp, h.Timestamp)
	buf = append(buf, tmp...)

	binary.LittleEndian.PutUint32(tmp, h.Bits)
	buf = append(buf, tmp...)

	binary.LittleEndian.PutUint32(tmp, h.Nonce)
	buf = append(buf, tmp...)

	return buf
}

func (h *BlockHeader) Hash() []byte {
	first := sha256.Sum256(h.Serialize())
	second := sha256.Sum256(first[:])
	return second[:]
}

func (h *BlockHeader) Validate() error {
	if len(h.PrevBlockHash) != 32 {
		return errors.New("invalid prev block hash length")
	}
	if len(h.MerkleRoot) != 32 {
		return errors.New("invalid merkle root length")
	}
	return nil
}

type Block struct {
	Header       BlockHeader
	Transactions []*Transaction
	Height       int
}

func (b *Block) Serialize() []byte {
	buf := bytes.Buffer{}
	buf.Write(b.Header.Serialize())

	buf.Write(encodeVarInt(uint64(len(b.Transactions))))
	for _, tx := range b.Transactions {
		buf.Write(tx.Serialize())
	}

	return buf.Bytes()
}

func (b *Block) BuildMerkleRoot() []byte {
	var txIDs [][]byte

	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.Serialize())
	}

	return BuildMerkleRoot(txIDs)
}

func NewBlock(txs []*Transaction, prevHash []byte, height int, bits uint32) *Block {
	block := &Block{
		Header: BlockHeader{
			Version:       1,
			PrevBlockHash: prevHash,
			MerkleRoot:    nil,
			Timestamp:     uint32(time.Now().Unix()),
			Bits:          bits,
			Nonce:         0,
		},
		Transactions: txs,
		Height:       height,
	}

	block.Header.MerkleRoot = block.BuildMerkleRoot()

	pow := NewProofOfWork(&block.Header)
	nonce, _ := pow.Run()

	block.Header.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *Transaction, bits uint32) *Block {
	return NewBlock([]*Transaction{coinbase}, make([]byte, 32), 0, bits)
}

func DeserializeBlockHeader(data []byte) (*BlockHeader, error) {
	if len(data) < 80 {
		return nil, errors.New("invalid block header length")
	}

	offset := 0
	h := &BlockHeader{}

	h.Version = binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	h.PrevBlockHash = make([]byte, 32)
	copy(h.PrevBlockHash, data[offset:offset+32])
	offset += 32

	h.MerkleRoot = make([]byte, 32)
	copy(h.MerkleRoot, data[offset:offset+32])
	offset += 32

	h.Timestamp = binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	h.Bits = binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	h.Nonce = binary.LittleEndian.Uint32(data[offset:])

	return h, nil
}

func DeserializeBlockHeaderFromReader(r io.Reader) (*BlockHeader, error) {
	h := &BlockHeader{
		PrevBlockHash: make([]byte, 32),
		MerkleRoot:    make([]byte, 32),
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(r, h.PrevBlockHash); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(r, h.MerkleRoot); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Timestamp); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Bits); err != nil {
		return nil, err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.Nonce); err != nil {
		return nil, err
	}

	return h, nil
}
