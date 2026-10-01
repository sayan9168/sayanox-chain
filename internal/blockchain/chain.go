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
	store  PersistentStore
	mu     sync.RWMutex
}

func NewChain() *Chain {
	return &Chain{
		Blocks: []*Block{NewGenesisBlock()},
		State:  state.NewState(),
	}
}

func NewPersistentChain(store PersistentStore) (*Chain, error) {
	blocks, err := store.LoadBlocks()
	if err != nil {
		return nil, err
	}

	chain := &Chain{
		Blocks: blocks,
		State:  state.NewState(),
		store:  store,
	}

	if len(chain.Blocks) == 0 {
		chain.Blocks = []*Block{NewGenesisBlock()}
		if err := store.Commit(chain.Blocks[0], chain.State); err != nil {
			return nil, err
		}
	} else {
		for i, block := range chain.Blocks {
			var previous *Block
			if i > 0 {
				previous = chain.Blocks[i-1]
			}
			if err := block.Validate(previous); err != nil {
				return nil, err
			}
		}
	}

	if err := store.LoadState(chain.State); err != nil {
		return nil, err
	}

	return chain, nil
}

func (c *Chain) AddTransaction(tx *transaction.Transaction) error {
	if err := transaction.Validate(tx); err != nil {
		return err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, block := range c.Blocks {
		for _, existing := range block.Transactions {
			if existing.Hash == tx.Hash {
				return errors.New("duplicate transaction")
			}
		}
	}
	return c.State.ValidateTransactions([]*transaction.Transaction{tx})
}

func (c *Chain) AddBlock(block *Block) error {
	if block == nil {
		return errors.New("nil block")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	previous := c.Blocks[len(c.Blocks)-1]
	if err := block.Validate(previous); err != nil {
		return err
	}

	nextState, err := c.State.Project(block.Transactions)
	if err != nil {
		return err
	}

	if c.store != nil {
		if err := c.store.Commit(block, nextState); err != nil {
			return err
		}
	}

	c.State = nextState
	c.Blocks = append(c.Blocks, block)
	return nil
}

func (c *Chain) LatestBlock() *Block {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Blocks) == 0 {
		return nil
	}
	return c.Blocks[len(c.Blocks)-1]
}

func (c *Chain) GetBlock(height uint64) (*Block, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if height >= uint64(len(c.Blocks)) {
		return nil, false
	}
	return c.Blocks[height], true
}

func (c *Chain) GetHeight() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Blocks)
}

func (c *Chain) Validate() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Blocks) == 0 {
		return errors.New("empty chain")
	}
	for i, block := range c.Blocks {
		var previous *Block
		if i > 0 {
			previous = c.Blocks[i-1]
		}
		if err := block.Validate(previous); err != nil {
			return err
		}
	}
	return nil
}
