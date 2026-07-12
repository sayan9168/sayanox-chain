package transaction

import "errors"

func Validate(tx *Transaction) error {

	if tx == nil {
		return errors.New("transaction is nil")
	}

	if tx.From == "" {
		return errors.New("missing sender address")
	}

	if tx.To == "" {
		return errors.New("missing receiver address")
	}

	if tx.Amount == nil || tx.Amount.Value.Sign() <= 0 {
		return errors.New("invalid transaction amount")
	}

	if tx.Fee == nil || tx.Fee.Value.Sign() < 0 {
		return errors.New("invalid transaction fee")
	}

	if tx.Signature == "" {
		return errors.New("missing signature")
	}

	if tx.PublicKey == "" {
		return errors.New("missing public key")
	}

	if !tx.Verify() {
		return errors.New("invalid signature")
	}

	return nil
}
