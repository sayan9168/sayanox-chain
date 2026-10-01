package transaction

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"

	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

func Validate(tx *Transaction) error {
	if tx == nil { return errors.New("transaction is nil") }
	if tx.From == "" || tx.To == "" { return errors.New("missing sender or receiver address") }
	if tx.Amount == nil || tx.Amount.Value == nil || tx.Amount.Value.Sign() <= 0 { return errors.New("invalid transaction amount") }
	if tx.Fee == nil || tx.Fee.Value == nil || tx.Fee.Value.Sign() < 0 { return errors.New("invalid transaction fee") }
	if tx.Timestamp <= 0 { return errors.New("invalid transaction timestamp") }
	if tx.Hash == "" || tx.Hash != tx.CalculateHash() { return errors.New("invalid transaction hash") }
	if tx.Signature == "" || tx.PublicKey == "" { return errors.New("missing signature") }

	pub, err := hex.DecodeString(tx.PublicKey)
	if err != nil || len(pub) != ed25519.PublicKeySize { return errors.New("invalid public key") }
	sig, err := hex.DecodeString(tx.Signature)
	if err != nil || len(sig) != ed25519.SignatureSize { return errors.New("invalid signature encoding") }
	if wallet.GenerateAddress(pub) != tx.From { return errors.New("sender does not match public key") }
	if !ed25519.Verify(ed25519.PublicKey(pub), []byte(tx.CalculateHash()), sig) { return errors.New("invalid signature") }
	return nil
}
