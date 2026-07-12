package core

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	genesis := NewBlock(0, "0", []string{"Genesis Block"})

	return &Blockchain{
		Blocks: []*Block{genesis},
	}
}

func (bc *Blockchain) AddBlock(transactions []string) {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := NewBlock(
		lastBlock.Index+1,
		lastBlock.Hash,
		transactions,
	)

	bc.Blocks = append(bc.Blocks, newBlock)
}
