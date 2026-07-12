package network

import (
	"errors"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
)

type SyncManager struct {
	Chain *blockchain.Chain
}

func NewSyncManager(chain *blockchain.Chain) *SyncManager {
	return &SyncManager{
		Chain: chain,
	}
}

func (s *SyncManager) GetHeight() int {
	return s.Chain.GetHeight()
}

func (s *SyncManager) CompareHeight(remoteHeight int) bool {

	localHeight := s.GetHeight()

	return remoteHeight > localHeight
}

func (s *SyncManager) ValidateSync() error {

	if s.Chain == nil {
		return errors.New("chain not initialized")
	}

	return nil
}
