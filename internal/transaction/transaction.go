package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Transaction struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    *Amount `json:"amount"`
	Fee       *Amount `json:"fee"`
	Timestamp int64   `json:"timestamp"`
	Hash      string  `json:"hash"`
	Signature string  `json:"signature"`
	PublicKey string  `json:"public_key"`
}

func NewTransaction(from, to string, amount, fee *Amount) *Transaction {
	tx := &Transaction{From: from, To: to, Amount: amount, Fee: fee, Timestamp: time.Now().UnixNano()}
	tx.Hash = tx.CalculateHash()
	return tx
}

func (tx *Transaction) CalculateHash() string {
	if tx == nil || tx.Amount == nil || tx.Fee == nil || tx.Amount.Value == nil || tx.Fee.Value == nil {
		return ""
	}
	data := fmt.Sprintf("%s|%s|%s|%s|%d", tx.From, tx.To, tx.Amount.String(), tx.Fee.String(), tx.Timestamp)
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}
