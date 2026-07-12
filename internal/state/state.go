package state

import (
	"errors"
	"math/big"
	"sync"

	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

type Account struct {
	Address string
	Balance *big.Int
}

type State struct {
	Accounts map[string]*Account
	mu       sync.RWMutex
}

func NewState() *State {
	return &State{
		Accounts: make(map[string]*Account),
	}
}

func (s *State) CreateAccount(address string, balance *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Accounts[address] = &Account{
		Address: address,
		Balance: new(big.Int).Set(balance),
	}
}

func (s *State) GetBalance(address string) *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	account, exists := s.Accounts[address]

	if !exists {
		return big.NewInt(0)
	}

	return new(big.Int).Set(account.Balance)
}

func (s *State) Transfer(tx *transaction.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sender, exists := s.Accounts[tx.From]

	if !exists {
		return errors.New("sender account not found")
	}

	totalCost := new(big.Int).Add(
		tx.Amount.Value,
		tx.Fee.Value,
	)

	if sender.Balance.Cmp(totalCost) < 0 {
		return errors.New("insufficient balance")
	}

	receiver, exists := s.Accounts[tx.To]

	if !exists {
		receiver = &Account{
			Address: tx.To,
			Balance: big.NewInt(0),
		}

		s.Accounts[tx.To] = receiver
	}

	sender.Balance.Sub(sender.Balance, totalCost)
	receiver.Balance.Add(receiver.Balance, tx.Amount.Value)

	return nil
}
