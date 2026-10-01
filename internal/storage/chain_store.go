package storage

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/state"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	stateKey  = "state"
	heightKey = "meta:height"
)

type ChainStore struct {
	db *LevelDB
}

func NewChainStore(db *LevelDB) *ChainStore {
	return &ChainStore{db: db}
}

func (s *ChainStore) LoadBlocks() ([]*blockchain.Block, error) {
	data, err := s.db.Get(heightKey)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var height uint64
	if err := json.Unmarshal(data, &height); err != nil {
		return nil, err
	}

	blocks := make([]*blockchain.Block, 0, height+1)
	for i := uint64(0); i <= height; i++ {
		block, err := s.GetBlock(i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (s *ChainStore) GetBlock(height uint64) (*blockchain.Block, error) {
	data, err := s.db.Get(fmt.Sprintf("block:%d", height))
	if err != nil {
		return nil, err
	}
	var block blockchain.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func (s *ChainStore) LoadState(stateData *state.State) error {
	data, err := s.db.Get(stateKey)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil
		}
		return err
	}
	var accounts map[string]*state.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return err
	}
	stateData.SetAccounts(accounts)
	return nil
}

func (s *ChainStore) Commit(block *blockchain.Block, stateData *state.State) error {
	blockData, err := json.Marshal(block)
	if err != nil {
		return err
	}
	accounts := stateData.SnapshotAccounts()
	stateDataBytes, err := json.Marshal(accounts)
	if err != nil {
		return err
	}
	heightData, err := json.Marshal(block.Height)
	if err != nil {
		return err
	}

	batch := new(leveldb.Batch)
	batch.Put([]byte(fmt.Sprintf("block:%d", block.Height)), blockData)
	batch.Put([]byte(heightKey), heightData)
	batch.Put([]byte(stateKey), stateDataBytes)

	return s.db.Write(batch)
}

func (s *ChainStore) Close() error {
	return s.db.Close()
}

func (s *ChainStore) Keys(prefix string) *util.Range {
	return util.BytesPrefix([]byte(prefix))
}
