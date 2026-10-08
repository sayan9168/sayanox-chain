package blockchain

import (
	"crypto/ed25519"
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
	Height             uint64                    `json:"height"`
	Timestamp          int64                     `json:"timestamp"`
	PreviousHash       string                    `json:"previous_hash"`
	MerkleRoot         string                    `json:"merkle_root"`
	Hash               string                    `json:"hash"`
	Proposer           string                    `json:"proposer,omitempty"`
	ProposerPublicKey  string                    `json:"proposer_public_key,omitempty"`
	ProposerSignature  string                    `json:"proposer_signature,omitempty"`
	Transactions       []*transaction.Transaction `json:"transactions"`
}

func NewBlock(height uint64, previousHash string, transactions []*transaction.Transaction) *Block {
	block := &Block{
		Height: height,
		Timestamp: time.Now().Unix(),
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
	block := &Block{Height: 0, Timestamp: timestamp.Unix(), Transactions: nil}
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

func (b *Block) signingPayload() ([]byte, error) {
	return json.Marshal(struct {
		Height       uint64
		Timestamp    int64
		PreviousHash string
		MerkleRoot   string
		Transactions []*transaction.Transaction
		Proposer     string
		ProposerPublicKey string
	}{
		b.Height, b.Timestamp, b.PreviousHash, b.MerkleRoot, b.Transactions,
		b.Proposer, b.ProposerPublicKey,
	})
}

func (b *Block) SigningHash() string {
	data, err := b.signingPayload()
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
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
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func (b *Block) Sign(privateKey ed25519.PrivateKey, address string) error {
	if b == nil || len(privateKey) != ed25519.PrivateKeySize {
		return errors.New("invalid block or private key")
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)
	derived, err := walletAddress(publicKey)
	if err != nil || derived != address {
		return errors.New("private key does not match proposer address")
	}
	b.Proposer = address
	b.ProposerPublicKey = hex.EncodeToString(publicKey)
	b.Hash = b.CalculateHash()
	b.ProposerSignature = hex.EncodeToString(ed25519.Sign(privateKey, []byte(b.SigningHash())))
	return nil
}

func walletAddress(publicKey ed25519.PublicKey) (string, error) {
	sum := sha256.Sum256(publicKey)
	return "sx1" + hex.EncodeToString(sum[:20]), nil
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
