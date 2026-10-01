package mempool

import (
	"sync"
	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

type Mempool struct { Transactions []*transaction.Transaction; mu sync.Mutex }

func NewMempool() *Mempool { return &Mempool{Transactions: make([]*transaction.Transaction,0)} }

func (m *Mempool) Add(tx *transaction.Transaction) {
	m.mu.Lock(); defer m.mu.Unlock()
	for _, existing := range m.Transactions { if existing.Hash == tx.Hash { return } }
	m.Transactions = append(m.Transactions, tx)
}

func (m *Mempool) GetAll() []*transaction.Transaction {
	m.mu.Lock(); defer m.mu.Unlock()
	return append([]*transaction.Transaction(nil), m.Transactions...)
}

func (m *Mempool) Clear() { m.mu.Lock(); defer m.mu.Unlock(); m.Transactions = nil }
