package blockchain

import (
	"math/big"

	"github.com/sayan9168/sayanox-chain/internal/state"
	"github.com/sayan9168/sayanox-chain/internal/transaction"
)

func CreateGenesisBlock(s *state.State, address string) *Block {

	initialSupply := new(big.Int)

	// 1,000,000,000 SAYX with 18 decimals
	initialSupply.SetString(
		"1000000000000000000000000000",
		10,
	)

	s.CreateAccount(
		address,
		initialSupply,
	)

	return NewBlock(
		0,
		"0",
		[]*transaction.Transaction{
			{
				From:      "GENESIS",
				To:        address,
				Amount:    &transaction.Amount{Value: initialSupply},
				Fee:       &transaction.Amount{Value: big.NewInt(0)},
				Timestamp: 0,
			},
		},
	)
}
