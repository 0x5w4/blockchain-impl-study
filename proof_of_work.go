package main

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/big"
)

type ProofOfWork struct {
	header *BlockHeader
	target *big.Int
}

func (pow *ProofOfWork) Run() (uint32, []byte) {
	var hashInt big.Int
	headerBytes := pow.header.Serialize()

	for nonce := range uint32(math.MaxUint32) {
		if len(headerBytes) != 80 {
			panic("Header must be exactly 80 bytes for Bitcoin-style PoW")
		}
		binary.LittleEndian.PutUint32(headerBytes[76:], nonce)

		first := sha256.Sum256(headerBytes)
		second := sha256.Sum256(first[:])

		hashInt.SetBytes(second[:])

		if hashInt.Cmp(pow.target) == -1 {
			return nonce, second[:]
		}
	}

	panic("nonce exhausted")
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	hash := pow.header.Hash()
	hashInt.SetBytes(hash)

	return hashInt.Cmp(pow.target) == -1
}

func NewProofOfWork(header *BlockHeader) *ProofOfWork {
	return &ProofOfWork{
		header: header,
		target: BitsToTarget(header.Bits),
	}
}

func BitsToTarget(bits uint32) *big.Int {
	exponent := bits >> 24
	mantissa := bits & 0xFFFFFF

	if mantissa == 0 || exponent < 3 {
		panic("invalid bits")
	}

	target := new(big.Int).SetUint64(uint64(mantissa))
	shift := 8 * (int(exponent) - 3)
	target.Lsh(target, uint(shift))

	return target
}
