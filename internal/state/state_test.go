package state

import (
	"math/big"
	"testing"

	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

func TestApplyTransactionsIsAtomic(t *testing.T) {
	s := NewState()
	s.CreateAccount("alice", big.NewInt(100))

	amount, _ := transaction.NewAmount("60")
	fee, _ := transaction.NewAmount("1")
	valid := transaction.NewTransaction("alice", "bob", amount, fee)
	valid.PublicKey = ""
	valid.Signature = ""
	valid.Hash = valid.CalculateHash()

	invalidAmount, _ := transaction.NewAmount("1000")
	invalid := transaction.NewTransaction("alice", "carol", invalidAmount, fee)
	invalid.Hash = invalid.CalculateHash()

	if err := s.ApplyTransactions([]*transaction.Transaction{valid, invalid}); err == nil {
		t.Fatal("expected atomic batch failure")
	}
	if s.GetBalance("alice").Cmp(big.NewInt(100)) != 0 {
		t.Fatal("state changed after failed batch")
	}
}
