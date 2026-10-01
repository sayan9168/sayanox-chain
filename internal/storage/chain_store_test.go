package storage

import (
	"math/big"
	"testing"

	"github.com/sayan9168/sayanox-chain/internal/state"
)

func TestChainStoreCommitAndLoad(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewChainStore(db)
	s := state.NewState()
	s.CreateAccount("alice", big.NewInt(42))

	blockData := []byte(`{"height":0,"timestamp":1,"previous_hash":"","merkle_root":"","hash":"","transactions":null}`)
	if err := db.Put("block:0", blockData); err != nil {
		t.Fatal(err)
	}
	if err := db.Put("meta:height", []byte("0")); err != nil {
		t.Fatal(err)
	}
	if err := db.Put("state", []byte(`{"alice":{"address":"alice","balance":42}}`)); err != nil {
		t.Fatal(err)
	}

	if _, err := store.LoadBlocks(); err != nil {
		t.Fatal(err)
	}
	var loaded state.State
	loaded = *state.NewState()
	if err := store.LoadState(&loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.GetBalance("alice").Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("unexpected balance: %s", loaded.GetBalance("alice"))
	}
}
