package transaction

import (
	"encoding/hex"
	"errors"

	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

func Validate(tx *Transaction) error {
	if tx == nil {
		return errors.New("nil transaction")
	}
	if tx.From == "" || tx.To == "" {
		return errors.New("missing sender or recipient")
	}
	if tx.From == tx.To {
		return errors.New("sender and recipient must differ")
	}
	if tx.Amount == nil || tx.Amount.Value == nil || tx.Amount.Value.Sign() <= 0 {
		return errors.New("invalid amount")
	}
	if tx.Fee == nil || tx.Fee.Value == nil || tx.Fee.Value.Sign() < 0 {
		return errors.New("invalid fee")
	}
	if tx.Timestamp <= 0 {
		return errors.New("invalid timestamp")
	}
	if tx.Hash == "" || tx.CalculateHash() != tx.Hash {
		return errors.New("invalid transaction hash")
	}
	if !tx.Verify() {
		return errors.New("invalid transaction signature")
	}
	address, err := wallet.AddressFromPublicKey(publicKeyBytes(tx.PublicKey))
	if err != nil || address != tx.From {
		return errors.New("public key does not match sender")
	}
	return nil
}

func publicKeyBytes(encoded string) []byte {
	b, err := hex.DecodeString(encoded)
	if err != nil {
		return nil
	}
	return b
}
