package consensus

import (
	"math/big"
)

type Validator struct {
	Address string
	Stake   *big.Int
}

type PoS struct {
	Validators []*Validator
}

func NewPoS() *PoS {
	return &PoS{
		Validators: make([]*Validator, 0),
	}
}

func (p *PoS) AddValidator(
	address string,
	stake *big.Int,
) {
	p.Validators = append(
		p.Validators,
		&Validator{
			Address: address,
			Stake:   new(big.Int).Set(stake),
		},
	)
}

func (p *PoS) TotalStake() *big.Int {

	total := big.NewInt(0)

	for _, validator := range p.Validators {
		total.Add(
			total,
			validator.Stake,
		)
	}

	return total
}
