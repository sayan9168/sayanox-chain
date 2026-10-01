package mempool

import (
	"errors"
	"sync"

	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

var (
	ErrNilTransaction   = errors.New("nil transaction")
	ErrDuplicate       = errors.New("duplicate transaction")
	ErrInvalidSignature = errors.New("invalid transaction")
)

type Mempool struct {
	Transactions []*transaction.Transaction
	byHash       map[string]struct{}
	mu           sync.RWMutex
}

func NewMempool() *Mempool {
	return &Mempool{
		Transactions: make([]*transaction.Transaction, 0),
		byHash:       make(map[string]struct{}),
	}
}

func (m *Mempool) Add(tx *transaction.Transaction) error {
	if tx == nil {
		return ErrNilTransaction
	}
	if err := transaction.Validate(tx); err != nil {
		return errors.Join(ErrInvalidSignature, err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.byHash[tx.Hash]; exists {
		return ErrDuplicate
	}
	m.Transactions = append(m.Transactions, tx)
	m.byHash[tx.Hash] = struct{}{}
	return nil
}

func (m *Mempool) GetAll() []*transaction.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*transaction.Transaction, len(m.Transactions))
	copy(out, m.Transactions)
	return out
}

func (m *Mempool) Remove(hash string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, tx := range m.Transactions {
		if tx.Hash != hash {
			continue
		}
		m.Transactions = append(m.Transactions[:i], m.Transactions[i+1:]...)
		delete(m.byHash, hash)
		return true
	}
	return false
}

func (m *Mempool) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Transactions = m.Transactions[:0]
	clear(m.byHash)
}

func (m *Mempool) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.Transactions)
}
