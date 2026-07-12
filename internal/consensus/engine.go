package consensus

import (
	"errors"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
)

type Engine struct {
	Name string
}

func NewEngine() *Engine {
	return &Engine{
		Name: "Proof-of-Stake",
	}
}

func (e *Engine) ValidateBlock(block *blockchain.Block) error {

	if block == nil {
		return errors.New("block is nil")
	}

	if block.Hash == "" {
		return errors.New("missing block hash")
	}

	return nil
}
