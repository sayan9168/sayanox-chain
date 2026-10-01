package blockchain

import "testing"

func TestGenesisBlock(t *testing.T) {
	chain := NewChain()
	if chain.GetHeight() != 1 {
		t.Fatalf("expected genesis height count 1, got %d", chain.GetHeight())
	}
	genesis := chain.LatestBlock()
	if genesis == nil {
		t.Fatal("missing genesis block")
	}
	if err := genesis.Validate(nil); err != nil {
		t.Fatalf("genesis validation failed: %v", err)
	}
	if err := chain.Validate(); err != nil {
		t.Fatalf("chain validation failed: %v", err)
	}
}

func TestBlockRejectsBrokenLink(t *testing.T) {
	chain := NewChain()
	block := NewBlock(1, "wrong", nil)
	if err := chain.AddBlock(block); err == nil {
		t.Fatal("expected broken previous hash to be rejected")
	}
}
