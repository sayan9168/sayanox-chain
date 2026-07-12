package blockchain

import (
	"errors"
	"sync"

	"github.com/sayan9168/sayanox-chain/internal/state"
	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

type Chain struct {
	Blocks []*Block
	State  *state.State
	mu     sync.Mutex
}

func NewChain() *Chain {
	return &Chain{
		Blocks: []*Block{},
		State:  state.NewState(),
	}
}

func (c *Chain) AddTransaction(tx *transaction.Transaction) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if tx.Amount.Value.Sign() <= 0 {
		return errors.New("invalid amount")
	}

	return c.State.Transfer(tx)
}

func (c *Chain) GetHeight() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.Blocks)
}
