package main

import (
	"crypto/sha256"
)

func hashPair(left, right []byte) []byte {
	first := sha256.Sum256(append(left, right...))
	second := sha256.Sum256(first[:])
	return second[:]
}

func BuildMerkleRoot(txHashes [][]byte) []byte {
	if len(txHashes) == 0 {
		return []byte{}
	}

	nodes := txHashes

	for len(nodes) > 1 {
		var level [][]byte

		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				level = append(level, hashPair(nodes[i], nodes[i+1]))
			} else {
				level = append(level, hashPair(nodes[i], nodes[i]))
			}
		}
		nodes = level
	}

	return nodes[0]
}
