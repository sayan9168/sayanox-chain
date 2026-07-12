package storage

import "github.com/sayan9168/sayanox-chain/core"

type BlockStore struct {
	db *Database
}

func NewBlockStore(db *Database) *BlockStore {
	return &BlockStore{db: db}
}

func (bs *BlockStore) SaveBlock(block *core.Block) {
	// TODO: Serialize and store the block.
}

func (bs *BlockStore) LoadBlock(hash string) (*core.Block, error) {
	// TODO: Load and deserialize the block.
	return nil, nil
}
