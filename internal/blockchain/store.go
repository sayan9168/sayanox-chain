package blockchain

import "github.com/sayan9168/sayanox-chain/internal/state"

type PersistentStore interface {
	LoadBlocks() ([]*Block, error)
	LoadState(*state.State) error
	Commit(*Block, *state.State) error
}
