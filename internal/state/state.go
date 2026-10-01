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
	Accounts map[string]*Account
	mu       sync.RWMutex
}

func NewState() *State {
	return &State{Accounts: make(map[string]*Account)}
}

func (s *State) CreateAccount(address string, balance *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Accounts[address] = &Account{Address: address, Balance: new(big.Int).Set(balance)}
}

func (s *State) GetBalance(address string) *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	account, exists := s.Accounts[address]
	if !exists || account.Balance == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(account.Balance)
}

func (s *State) ValidateTransfer(tx *transaction.Transaction) error {
	if tx == nil || tx.Amount == nil || tx.Fee == nil {
		return errors.New("invalid transaction")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	sender, exists := s.Accounts[tx.From]
	if !exists {
		return errors.New("sender account not found")
	}
	total := new(big.Int).Add(tx.Amount.Value, tx.Fee.Value)
	if sender.Balance.Cmp(total) < 0 {
		return errors.New("insufficient balance")
	}
	return nil
}

func (s *State) Transfer(tx *transaction.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return transferLocked(s.Accounts, tx)
}

func (s *State) ApplyTransactions(txs []*transaction.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := cloneAccounts(s.Accounts)
	for _, tx := range txs {
		if err := transferLocked(next, tx); err != nil {
			return err
		}
	}
	s.Accounts = next
	return nil
}

func cloneAccounts(src map[string]*Account) map[string]*Account {
	dst := make(map[string]*Account, len(src))
	for address, account := range src {
		dst[address] = &Account{
			Address: account.Address,
			Balance: new(big.Int).Set(account.Balance),
		}
	}
	return dst
}

func transferLocked(accounts map[string]*Account, tx *transaction.Transaction) error {
	if tx == nil || tx.Amount == nil || tx.Fee == nil {
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
	receiver, exists := accounts[tx.To]
	if !exists {
		receiver = &Account{Address: tx.To, Balance: big.NewInt(0)}
		accounts[tx.To] = receiver
	}
	sender.Balance.Sub(sender.Balance, total)
	receiver.Balance.Add(receiver.Balance, tx.Amount.Value)
	return nil
}
