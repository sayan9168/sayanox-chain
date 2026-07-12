package consensus

import (
	"errors"
	"math/big"
)

func (p *PoS) SelectValidator(seed *big.Int) (*Validator, error) {

	if len(p.Validators) == 0 {
		return nil, errors.New("no validators available")
	}

	totalStake := p.TotalStake()

	if totalStake.Sign() == 0 {
		return nil, errors.New("total stake is zero")
	}

	target := new(big.Int).Mod(seed, totalStake)

	current := big.NewInt(0)

	for _, validator := range p.Validators {

		current.Add(
			current,
			validator.Stake,
		)

		if current.Cmp(target) > 0 {
			return validator, nil
		}
	}

	return p.Validators[len(p.Validators)-1], nil
}
