package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/sayan9168/sayanox-chain/internal/state"
)

type StateStore struct {
	FilePath string
	mu       sync.Mutex
}

func NewStateStore(path string) *StateStore {
	return &StateStore{
		FilePath: path,
	}
}

func (s *StateStore) Save(stateData *state.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(stateData.Accounts)
	if err != nil {
		return err
	}

	return os.WriteFile(s.FilePath, data, 0644)
}

func (s *StateStore) Load(stateData *state.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &stateData.Accounts)
}
