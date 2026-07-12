package transaction

import (
	"math/big"
)

type Amount struct {
	Value *big.Int
}

func NewAmount(value string) (*Amount, bool) {
	number, ok := new(big.Int).SetString(value, 10)

	if !ok {
		return nil, false
	}

	return &Amount{
		Value: number,
	}, true
}

func (a *Amount) Add(other *Amount) *Amount {
	return &Amount{
		Value: new(big.Int).Add(a.Value, other.Value),
	}
}

func (a *Amount) Sub(other *Amount) *Amount {
	return &Amount{
		Value: new(big.Int).Sub(a.Value, other.Value),
	}
}

func (a *Amount) String() string {
	return a.Value.String()
}
