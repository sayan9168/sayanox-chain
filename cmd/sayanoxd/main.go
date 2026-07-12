package main

import (
	"fmt"

	"github.com/sayan9168/sayanox-chain/core"
)

func main() {
	// Create blockchain
	bc := core.NewBlockchain()

	// Add sample blocks
	bc.AddBlock([]string{"Alice -> Bob : 10 SAYX"})
	bc.AddBlock([]string{"Bob -> Charlie : 5 SAYX"})

	fmt.Println("=== Sayanox Chain ===")

	for _, block := range bc.Blocks {
		fmt.Println("---------------------------")
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Transactions: %v\n", block.Transactions)
	}
}
