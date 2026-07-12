package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        uint64
	Timestamp    int64
	PreviousHash string
	Hash         string
	Transactions []string
}

func NewBlock(index uint64, previousHash string, transactions []string) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		PreviousHash: previousHash,
		Transactions: transactions,
	}

	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) CalculateHash() string {
	data := fmt.Sprintf("%d%d%s%v",
		b.Index,
		b.Timestamp,
		b.PreviousHash,
		b.Transactions,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
