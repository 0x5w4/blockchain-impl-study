// package main

// import (
// 	"fmt"
// 	"math/big"
// 	"testing"
// )

// func TestManualMerkleAndMining(t *testing.T) {
// 	fmt.Println("--- Starting Integration Test ---")
// //
// 	// 1. Create Dummy Transactions
// 	// We manually construct transactions to test the Merkle Root calculation
// 	txIn := TXInput{[]byte{}, -1, nil, []byte("Coinbase Data")}
// 	txOut := TXOutput{50, []byte("Receiver Public Key")}
// 	coinbaseTx := &Transaction{
// 		Version:  1,
// 		ID:       []byte{}, // ID will be empty for this test, usually it's a hash
// 		Vin:      []TXInput{txIn},
// 		Vout:     []TXOutput{txOut},
// 		LockTime: 0,
// 	}
// 	// Manually set ID for the merkle tree to use (in a real app, this is the hash of the tx)
// 	coinbaseTx.ID = []byte("tx1_hash_32_bytes_long_123456789")

// 	// Create a second transaction
// 	tx2 := &Transaction{Version: 1, ID: []byte("tx2_hash_32_bytes_long_123456789"), Vin: nil, Vout: nil}

// 	transactions := []*Transaction{coinbaseTx, tx2}

// 	// 2. Verify Merkle Root Construction
// 	// We expect specific behavior: hashing pairs of IDs.
// 	// You can double-check this manual calculation if you want absolute certainty.
// 	merkleRoot := BuildMerkleRoot([][]byte{coinbaseTx.ID, tx2.ID})

// 	fmt.Printf("Merkle Root: %x\n", merkleRoot)

// 	if len(merkleRoot) != 32 {
// 		t.Fatalf("Error: Merkle Root must be 32 bytes, got %d", len(merkleRoot))
// 	}

// 	// 3. Create a Block
// 	// Difficulty "bits" - let's use a simpler target for the test so it mines fast
// 	// 0x1e00ffff is a standard "easy" difficulty for testing (requires somewhat low hash)
// 	// 0x207fffff would be the easiest possible (target is very high)
// 	bits := uint32(0x1d00ffff)

// 	fmt.Println("Mining block...")
// 	block := NewBlock(transactions, make([]byte, 32), 1, bits)

// 	// 4. Verify Header Length
// 	serializedHeader := block.Header.Serialize()
// 	if len(serializedHeader) != 80 {
// 		t.Fatalf("Error: Header length is %d, expected 80", len(serializedHeader))
// 	}
// 	fmt.Println("Header Length Verified: 80 bytes")

// 	// 5. Verify Proof of Work
// 	pow := NewProofOfWork(&block.Header)
// 	isValid := pow.Validate()

// 	fmt.Printf("PoW Valid: %v\n", isValid)
// 	fmt.Printf("Block Hash: %x\n", block.Header.Hash())
// 	fmt.Printf("Block Nonce: %d\n", block.Header.Nonce)

// 	if !isValid {
// 		t.Fatal("Error: Block failed Proof of Work validation!")
// 	}

// 	// 6. Verify Difficulty Target
// 	// Ensure the hash is actually below the target
// 	hashInt := new(big.Int)
// 	hashInt.SetBytes(block.Header.Hash())

// 	target := BitsToTarget(block.Header.Bits)

// 	if hashInt.Cmp(target) != -1 {
// 		t.Fatal("Error: Hash is NOT less than target (Mining logic broken)")
// 	}

// 	fmt.Println("--- Test PASSED Successfully ---")
// }

package main

import (
	"fmt"
	"os"
	"testing"
)

func TestBlockchainPersistence(t *testing.T) {
	nodeID := "test_node"
	os.RemoveAll("./tmp/blocks_" + nodeID) // Clean up

	// 1. Create Chain
	fmt.Println("Initializing Blockchain...")
	bc := InitBlockchain("test_address", nodeID)
	bc.Close()

	// 2. Re-open Chain
	fmt.Println("Re-opening Blockchain...")
	bc2 := ContinueBlockchain(nodeID)

	// 3. Add a Block
	fmt.Println("Mining new block...")
	tx := NewCoinbaseTX("test_address", "Block 2 Data")
	bc2.AddBlock([]*Transaction{tx})

	// 4. Verify Tip
	if len(bc2.LastHash) == 0 {
		t.Error("LastHash is empty")
	}
	fmt.Printf("Current Tip Hash: %x\n", bc2.LastHash)

	bc2.Close()

	// Clean up
	os.RemoveAll("./tmp/blocks_" + nodeID)
	fmt.Println("Persistence Test Passed!")
}
