package consensus

import (
	"crypto/ed25519"
	"errors"
	"math/big"
	"sort"
)

type Validator struct {
	Address   string
	Stake     *big.Int
	PublicKey ed25519.PublicKey
}

type PoS struct {
	Validators []*Validator
}

func NewPoS() *PoS {
	return &PoS{Validators: make([]*Validator, 0)}
}

func (p *PoS) AddValidator(address string, stake *big.Int) {
	p.AddValidatorWithKey(address, stake, nil)
}

func (p *PoS) AddValidatorWithKey(address string, stake *big.Int, publicKey ed25519.PublicKey) {
	if p == nil || address == "" || stake == nil || stake.Sign() <= 0 {
		return
	}
	key := append(ed25519.PublicKey(nil), publicKey...)
	for _, v := range p.Validators {
		if v.Address == address {
			v.Stake = new(big.Int).Set(stake)
			v.PublicKey = key
			return
		}
	}
	p.Validators = append(p.Validators, &Validator{
		Address: address,
		Stake: new(big.Int).Set(stake),
		PublicKey: key,
	})
	sort.Slice(p.Validators, func(i, j int) bool {
		return p.Validators[i].Address < p.Validators[j].Address
	})
}

func (p *PoS) TotalStake() *big.Int {
	total := new(big.Int)
	if p == nil {
		return total
	}
	for _, validator := range p.Validators {
		if validator != nil && validator.Stake != nil && validator.Stake.Sign() > 0 {
			total.Add(total, validator.Stake)
		}
	}
	return total
}

func (p *PoS) Validator(address string) (*Validator, bool) {
	if p == nil {
		return nil, false
	}
	for _, v := range p.Validators {
		if v != nil && v.Address == address {
			return v, true
		}
	}
	return nil, false
}

func (p *PoS) SelectValidator(seed *big.Int) (*Validator, error) {
	if p == nil || len(p.Validators) == 0 {
		return nil, errors.New("no validators available")
	}
	if seed == nil {
		return nil, errors.New("nil proposer seed")
	}
	total := p.TotalStake()
	if total.Sign() <= 0 {
		return nil, errors.New("total stake is zero")
	}
	target := new(big.Int).Mod(new(big.Int).Set(seed), total)
	current := new(big.Int)
	for _, validator := range p.Validators {
		if validator == nil || validator.Stake == nil || validator.Stake.Sign() <= 0 {
			continue
		}
		current.Add(current, validator.Stake)
		if target.Cmp(current) < 0 {
			return validator, nil
		}
	}
	return nil, errors.New("validator selection failed")
}
