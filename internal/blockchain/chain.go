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
	mu sync.RWMutex
}

func NewChain() *Chain {
	genesis := NewBlock(0, "0", nil)
	return &Chain{Blocks: []*Block{genesis}, State: state.NewState()}
}

func (c *Chain) AddTransaction(tx *transaction.Transaction) error {
	if err := transaction.Validate(tx); err != nil { return err }
	c.mu.RLock(); defer c.mu.RUnlock()
	for _, b := range c.Blocks {
		for _, existing := range b.Transactions {
			if existing.Hash == tx.Hash { return errors.New("duplicate transaction") }
		}
	}
	return c.State.ValidateTransactions([]*transaction.Transaction{tx})
}

func (c *Chain) AppendBlock(block *Block) error {
	c.mu.Lock(); defer c.mu.Unlock()
	var previous *Block
	if len(c.Blocks) > 0 { previous = c.Blocks[len(c.Blocks)-1] }
	if err := ValidateBlock(block, previous); err != nil { return err }
	if err := c.State.ApplyTransactions(block.Transactions); err != nil { return err }
	c.Blocks = append(c.Blocks, block)
	return nil
}

func (c *Chain) GetHeight() int {
	c.mu.RLock(); defer c.mu.RUnlock()
	if len(c.Blocks) == 0 { return 0 }
	return len(c.Blocks)-1
}

func (c *Chain) LatestBlock() *Block {
	c.mu.RLock(); defer c.mu.RUnlock()
	if len(c.Blocks)==0 { return nil }
	return c.Blocks[len(c.Blocks)-1]
}

func (c *Chain) GetBlock(height uint64) *Block {
	c.mu.RLock(); defer c.mu.RUnlock()
	if height >= uint64(len(c.Blocks)) { return nil }
	return c.Blocks[height]
}
