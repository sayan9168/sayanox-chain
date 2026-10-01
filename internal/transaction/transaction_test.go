package transaction

import (
	"testing"

	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

func TestTransactionSignAndVerify(t *testing.T) {
	w, err := wallet.NewWallet()
	if err != nil {
		t.Fatal(err)
	}
	amount, ok := NewAmount("1000000000000000000")
	if !ok {
		t.Fatal("failed to create amount")
	}
	fee, ok := NewAmount("1000000000000000")
	if !ok {
		t.Fatal("failed to create fee")
	}

	tx, err := NewSignedTransaction(w, "sx1recipient", amount, fee, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !tx.Verify() {
		t.Fatal("signature verification failed")
	}
	if err := Validate(tx); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	tx.Signature = "00"
	if tx.Verify() {
		t.Fatal("tampered signature accepted")
	}
}

func TestHashChangesWithNonce(t *testing.T) {
	w, err := wallet.NewWallet()
	if err != nil {
		t.Fatal(err)
	}
	amount, _ := NewAmount("1")
	fee, _ := NewAmount("0")
	tx1, _ := NewSignedTransaction(w, "sx1recipient", amount, fee, 1)
	tx2, _ := NewSignedTransaction(w, "sx1recipient", amount, fee, 2)
	if tx1.Hash == tx2.Hash {
		t.Fatal("nonce must affect transaction hash")
	}
}
