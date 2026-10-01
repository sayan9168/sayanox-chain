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
	mu       sync.RWMutex
}

func NewState() *State {
	return &State{Accounts: make(map[string]*Account)}
}

func (s *State) CreateAccount(address string, balance *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if balance == nil {
		balance = new(big.Int)
	}
	s.Accounts[address] = &Account{
		Address: address,
		Balance: new(big.Int).Set(balance),
	}
}

func (s *State) GetBalance(address string) *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	account, exists := s.Accounts[address]
	if !exists || account.Balance == nil {
		return new(big.Int)
	}
	return new(big.Int).Set(account.Balance)
}

func (s *State) SetAccounts(accounts map[string]*Account) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Accounts = cloneAccounts(accounts)
}

func (s *State) SnapshotAccounts() map[string]*Account {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneAccounts(s.Accounts)
}

func (s *State) ValidateTransactions(txs []*transaction.Transaction) error {
	_, err := s.Project(txs)
	return err
}

func (s *State) Project(txs []*transaction.Transaction) (*State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accounts := cloneAccounts(s.Accounts)
	for _, tx := range txs {
		if err := transaction.Validate(tx); err != nil {
			return nil, err
		}
		if err := validateTransfer(tx, accounts); err != nil {
			return nil, err
		}
		applyTransfer(accounts, tx)
	}

	return &State{Accounts: accounts}, nil
}

func (s *State) ApplyTransactions(txs []*transaction.Transaction) error {
	next, err := s.Project(txs)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Accounts = next.Accounts
	return nil
}

func (s *State) Transfer(tx *transaction.Transaction) error {
	return s.ApplyTransactions([]*transaction.Transaction{tx})
}

func cloneAccounts(src map[string]*Account) map[string]*Account {
	dst := make(map[string]*Account, len(src))
	for address, account := range src {
		balance := new(big.Int)
		if account != nil && account.Balance != nil {
			balance.Set(account.Balance)
		}
		dst[address] = &Account{
			Address: address,
			Balance: balance,
		}
	}
	return dst
}

func validateTransfer(tx *transaction.Transaction, accounts map[string]*Account) error {
	if tx == nil || tx.Amount == nil || tx.Fee == nil ||
		tx.Amount.Value == nil || tx.Fee.Value == nil {
		return errors.New("invalid transaction")
	}
	sender, exists := accounts[tx.From]
	if !exists {
		return errors.New("sender account not found")
	}
	total := new(big.Int).Add(tx.Amount.Value, tx.Fee.Value)
	if sender.Balance.Cmp(total) < 0 {
		return errors.New("insufficient balance")
	}
	return nil
}

func applyTransfer(accounts map[string]*Account, tx *transaction.Transaction) {
	total := new(big.Int).Add(tx.Amount.Value, tx.Fee.Value)
	sender := accounts[tx.From]
	receiver, exists := accounts[tx.To]
	if !exists {
		receiver = &Account{Address: tx.To, Balance: new(big.Int)}
		accounts[tx.To] = receiver
	}
	sender.Balance.Sub(sender.Balance, total)
	receiver.Balance.Add(receiver.Balance, tx.Amount.Value)
}
