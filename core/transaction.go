package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Transaction struct {
	From      string
	To        string
	Amount    uint64
	Fee       uint64
	Timestamp int64
	Hash      string
}

func NewTransaction(from, to string, amount, fee uint64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().Unix(),
	}

	tx.Hash = tx.CalculateHash()
	return tx
}

func (tx *Transaction) CalculateHash() string {
	data := fmt.Sprintf("%s%s%d%d%d",
		tx.From,
		tx.To,
		tx.Amount,
		tx.Fee,
		tx.Timestamp,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
