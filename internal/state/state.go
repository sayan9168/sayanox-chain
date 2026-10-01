package state

import (
	"errors"
	"math/big"
	"sync"

	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

type Account struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
}

type State struct {
	Accounts map[string]*Account `json:"accounts"`
	mu sync.RWMutex
}

func NewState() *State { return &State{Accounts: make(map[string]*Account)} }

func (s *State) CreateAccount(address string, balance *big.Int) {
	s.mu.Lock(); defer s.mu.Unlock()
	if balance == nil { balance = new(big.Int) }
	s.Accounts[address] = &Account{Address: address, Balance: new(big.Int).Set(balance)}
}

func (s *State) GetBalance(address string) *big.Int {
	s.mu.RLock(); defer s.mu.RUnlock()
	if a, ok := s.Accounts[address]; ok && a.Balance != nil { return new(big.Int).Set(a.Balance) }
	return new(big.Int)
}

func (s *State) validateTransfer(tx *transaction.Transaction, accounts map[string]*Account) error {
	sender, ok := accounts[tx.From]
	if !ok { return errors.New("sender account not found") }
	total := new(big.Int).Add(tx.Amount.Value, tx.Fee.Value)
	if sender.Balance.Cmp(total) < 0 { return errors.New("insufficient balance") }
	return nil
}

func cloneAccounts(src map[string]*Account) map[string]*Account {
	dst := make(map[string]*Account, len(src))
	for k, v := range src { dst[k] = &Account{Address: v.Address, Balance: new(big.Int).Set(v.Balance)} }
	return dst
}

func (s *State) ValidateTransactions(txs []*transaction.Transaction) error {
	s.mu.RLock(); defer s.mu.RUnlock()
	sim := cloneAccounts(s.Accounts)
	for _, tx := range txs {
		if err := transaction.Validate(tx); err != nil { return err }
		if err := s.validateTransfer(tx, sim); err != nil { return err }
		total := new(big.Int).Add(tx.Amount.Value, tx.Fee.Value)
		sim[tx.From].Balance.Sub(sim[tx.From].Balance, total)
		if _, ok := sim[tx.To]; !ok { sim[tx.To] = &Account{Address: tx.To, Balance: new(big.Int)} }
		sim[tx.To].Balance.Add(sim[tx.To].Balance, tx.Amount.Value)
	}
	return nil
}

func (s *State) ApplyTransactions(txs []*transaction.Transaction) error {
	s.mu.Lock(); defer s.mu.Unlock()
	sim := cloneAccounts(s.Accounts)
	for _, tx := range txs {
		if err := transaction.Validate(tx); err != nil { return err }
		if err := s.validateTransfer(tx, sim); err != nil { return err }
		total := new(big.Int).Add(tx.Amount.Value, tx.Fee.Value)
		sim[tx.From].Balance.Sub(sim[tx.From].Balance, total)
		if _, ok := sim[tx.To]; !ok { sim[tx.To] = &Account{Address: tx.To, Balance: new(big.Int)} }
		sim[tx.To].Balance.Add(sim[tx.To].Balance, tx.Amount.Value)
	}
	s.Accounts = sim
	return nil
}

func (s *State) Transfer(tx *transaction.Transaction) error { return s.ApplyTransactions([]*transaction.Transaction{tx}) }
