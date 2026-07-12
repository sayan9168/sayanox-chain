package core

type Validator struct {
	Address string
	Stake   uint64
	Active  bool
}

type ValidatorSet struct {
	Validators []*Validator
}

func NewValidatorSet() *ValidatorSet {
	return &ValidatorSet{
		Validators: []*Validator{},
	}
}

func (vs *ValidatorSet) AddValidator(address string, stake uint64) {
	validator := &Validator{
		Address: address,
		Stake:   stake,
		Active:  true,
	}

	vs.Validators = append(vs.Validators, validator)
}

func (vs *ValidatorSet) TotalStake() uint64 {
	var total uint64

	for _, v := range vs.Validators {
		if v.Active {
			total += v.Stake
		}
	}

	return total
}
