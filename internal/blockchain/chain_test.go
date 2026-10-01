package blockchain

import (
	"math/big"
	"testing"
	"github.com/sayan9168/sayanox-chain/internal/transaction"
	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

func TestChainAppendBlockAppliesStateAtomically(t *testing.T) {
	c := NewChain()
	w, err := wallet.NewWallet(); if err != nil { t.Fatal(err) }
	receiver, err := wallet.NewWallet(); if err != nil { t.Fatal(err) }
	c.State.CreateAccount(w.Address(), big.NewInt(100))

	amt, _ := transaction.NewAmount("10"); fee, _ := transaction.NewAmount("1")
	tx := transaction.NewTransaction(w.Address(), receiver.Address(), amt, fee)
	tx.Sign(w.PrivateKey)

	p := NewBlock(1, c.LatestBlock().Hash, []*transaction.Transaction{tx})
	if err := c.AppendBlock(p); err != nil { t.Fatal(err) }
	if c.GetHeight() != 1 { t.Fatalf("height=%d", c.GetHeight()) }
	if c.State.GetBalance(w.Address()).Cmp(big.NewInt(89)) != 0 { t.Fatal("sender balance incorrect") }
	if c.State.GetBalance(receiver.Address()).Cmp(big.NewInt(10)) != 0 { t.Fatal("receiver balance incorrect") }
}
