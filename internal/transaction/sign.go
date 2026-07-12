package transaction

import (
	"crypto/ed25519"
	"encoding/hex"
)

func (tx *Transaction) Sign(privateKey ed25519.PrivateKey) {

	hash := tx.CalculateHash()

	signature := ed25519.Sign(
		privateKey,
		[]byte(hash),
	)

	publicKey := privateKey.Public().(ed25519.PublicKey)

	tx.Signature = hex.EncodeToString(signature)
	tx.PublicKey = hex.EncodeToString(publicKey)
}

func (tx *Transaction) Verify() bool {

	signature, err := hex.DecodeString(tx.Signature)

	if err != nil {
		return false
	}

	publicKey, err := hex.DecodeString(tx.PublicKey)

	if err != nil {
		return false
	}

	return ed25519.Verify(
		ed25519.PublicKey(publicKey),
		[]byte(tx.CalculateHash()),
		signature,
	)
}
