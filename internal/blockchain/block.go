package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/sayan9168/sayanox-chain/genesis"
	"github.com/sayan9168/sayanox-chain/internal/merkle"
	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

type Block struct {
	Height       uint64                   `json:"height"`
	Timestamp    int64                    `json:"timestamp"`
	PreviousHash string                   `json:"previous_hash"`
	MerkleRoot   string                   `json:"merkle_root"`
	Hash         string                   `json:"hash"`
	Transactions []*transaction.Transaction `json:"transactions"`
}

func NewBlock(height uint64, previousHash string, transactions []*transaction.Transaction) *Block {
	block := &Block{
		Height:       height,
		Timestamp:    time.Now().Unix(),
		PreviousHash: previousHash,
		Transactions: append([]*transaction.Transaction(nil), transactions...),
	}
	block.MerkleRoot = merkle.Root(transactionHashes(transactions))
	block.Hash = block.CalculateHash()
	return block
}

func NewGenesisBlock() *Block {
	g := genesis.MustEmbedded()
	timestamp, err := time.Parse(time.RFC3339, g.GenesisTime)
	if err != nil {
		panic(err)
	}
	block := &Block{
		Height:       0,
		Timestamp:    timestamp.Unix(),
		Transactions: nil,
	}
	block.MerkleRoot = merkle.Root(nil)
	block.Hash = block.CalculateHash()
	return block
}

func transactionHashes(transactions []*transaction.Transaction) []string {
	leaves := make([]string, 0, len(transactions))
	for _, tx := range transactions {
		if tx != nil {
			leaves = append(leaves, tx.Hash)
		}
	}
	return leaves
}

func (b *Block) CalculateHash() string {
	data, err := json.Marshal(struct {
		Height       uint64
		Timestamp    int64
		PreviousHash string
		MerkleRoot   string
		Transactions []*transaction.Transaction
	}{
		b.Height, b.Timestamp, b.PreviousHash, b.MerkleRoot, b.Transactions,
	})
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (b *Block) Validate(previous *Block) error {
	if b == nil {
		return errors.New("nil block")
	}
	if b.Hash == "" || b.CalculateHash() != b.Hash {
		return errors.New("invalid block hash")
	}
	if merkle.Root(transactionHashes(b.Transactions)) != b.MerkleRoot {
		return errors.New("invalid merkle root")
	}
	if previous == nil {
		if b.Height != 0 || b.PreviousHash != "" {
			return errors.New("invalid genesis block")
		}
		return nil
	}
	if b.Height != previous.Height+1 {
		return errors.New("invalid block height")
	}
	if b.PreviousHash != previous.Hash {
		return errors.New("invalid previous hash")
	}
	if b.Timestamp < previous.Timestamp {
		return errors.New("block timestamp before previous block")
	}
	for _, tx := range b.Transactions {
		if err := transaction.Validate(tx); err != nil {
			return err
		}
	}
	return nil
}
