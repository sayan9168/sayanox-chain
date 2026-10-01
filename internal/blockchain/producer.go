package blockchain

import (
	"errors"

	"github.com/sayan9168/sayanox-chain/internal/mempool"
)

type Producer struct {
	Chain *Chain
	Mempool *mempool.Mempool
}

func NewProducer(chain *Chain, mp *mempool.Mempool) *Producer {
	return &Producer{Chain: chain, Mempool: mp}
}

func (p *Producer) CreateBlock() (*Block, error) {
	transactions := p.Mempool.GetAll()
	if len(transactions) == 0 { return nil, errors.New("no transactions available") }
	last := p.Chain.LatestBlock()
	height := uint64(0)
	previousHash := "0"
	if last != nil { height = last.Height + 1; previousHash = last.Hash }
	block := NewBlock(height, previousHash, transactions)
	if err := p.Chain.AppendBlock(block); err != nil { return nil, err }
	p.Mempool.Clear()
	return block, nil
}
