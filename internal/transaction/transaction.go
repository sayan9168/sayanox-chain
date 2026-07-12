package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Transaction struct {
	From      string
	To        string
	Amount    *Amount
	Fee       *Amount
	Timestamp int64
	Hash      string

        Signature string
        PublicKey string
}

func NewTransaction(from, to string, amount, fee *Amount) *Transaction {
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
	data := fmt.Sprintf(
		"%s%s%s%s%d",
		tx.From,
		tx.To,
		tx.Amount.String(),
		tx.Fee.String(),
		tx.Timestamp,
	)

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}
