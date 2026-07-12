package storage

import (
	"errors"
	"sync"
)

type Database struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string][]byte),
	}
}

func (db *Database) Put(key string, value []byte) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = value
}

func (db *Database) Get(key string) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	value, ok := db.data[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	return value, nil
}

func (db *Database) Delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, key)
}
