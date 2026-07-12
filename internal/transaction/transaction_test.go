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

	tx := NewTransaction(
		"sender",
		"receiver",
		amount,
		fee,
	)

	tx.Sign(w.PrivateKey)

	if !tx.Verify() {
		t.Fatal("signature verification failed")
	}

	if err := Validate(tx); err != nil {
		t.Fatalf("validation failed: %v", err)
	}
}
