package consensus

import (
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

func TestPoSSelectionIsDeterministic(t *testing.T) {
	p := NewPoS()
	p.AddValidator("sx1a", big.NewInt(10))
	p.AddValidator("sx1b", big.NewInt(20))

	a, err := p.SelectValidator(big.NewInt(7))
	if err != nil { t.Fatal(err) }
	b, err := p.SelectValidator(big.NewInt(7))
	if err != nil { t.Fatal(err) }
	if a.Address != b.Address { t.Fatalf("selection is not deterministic") }
}

func TestBlockProposerSignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil { t.Fatal(err) }
	w := &wallet.Wallet{PublicKey: pub, PrivateKey: priv}
	previous := blockchain.NewGenesisBlock()
	engine := NewEngine()
	engine.PoS.AddValidatorWithKey(w.Address(), big.NewInt(100), pub)

	block := blockchain.NewBlock(1, previous.Hash, nil)
	// Force the deterministic validator for this test by using the actual proposer check.
	proposer, err := engine.Proposer(previous)
	if err != nil { t.Fatal(err) }
	if proposer.Address != w.Address() { t.Fatal("unexpected proposer") }
	if err := block.Sign(priv, w.Address()); err != nil { t.Fatal(err) }
	if err := engine.ValidateBlock(block, previous); err != nil { t.Fatal(err) }
}
