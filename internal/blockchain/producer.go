package blockchain

import (
	"errors"

	"github.com/sayan9168/sayanox-chain/internal/mempool"
)

type Producer struct {
	Chain   *Chain
	Mempool *mempool.Mempool
}

func NewProducer(
	chain *Chain,
	mp *mempool.Mempool,
) *Producer {
	return &Producer{
		Chain:   chain,
		Mempool: mp,
	}
}

func (p *Producer) CreateBlock() (*Block, error) {

	transactions := p.Mempool.GetAll()

	if len(transactions) == 0 {
		return nil, errors.New("no transactions available")
	}

	var previousHash string
	height := uint64(len(p.Chain.Blocks))

	if height > 0 {
		previousHash = p.Chain.Blocks[height-1].Hash
	}

	block := NewBlock(
		height,
		previousHash,
		transactions,
	)

	err := ValidateBlock(
		block,
		p.getLastBlock(),
	)

	if err != nil {
		return nil, err
	}

	p.Chain.Blocks = append(
		p.Chain.Blocks,
		block,
	)

	p.Mempool.Clear()

	return block, nil
}

func (p *Producer) getLastBlock() *Block {

	if len(p.Chain.Blocks) == 0 {
		return nil
	}

	return p.Chain.Blocks[len(p.Chain.Blocks)-1]
}
