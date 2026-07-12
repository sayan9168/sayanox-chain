package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

type Block struct {
	Height       uint64
	Timestamp    int64
	PreviousHash string
	Hash         string
	Transactions []*transaction.Transaction
}

func NewBlock(
	height uint64,
	previousHash string,
	transactions []*transaction.Transaction,
) *Block {

	block := &Block{
		Height:       height,
		Timestamp:    time.Now().Unix(),
		PreviousHash: previousHash,
		Transactions: transactions,
	}

	block.Hash = block.CalculateHash()

	return block
}

func (b *Block) CalculateHash() string {
	data, _ := json.Marshal(struct {
		Height       uint64
		Timestamp    int64
		PreviousHash string
		Transactions []*transaction.Transaction
	}{
		b.Height,
		b.Timestamp,
		b.PreviousHash,
		b.Transactions,
	})

	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:])
}
