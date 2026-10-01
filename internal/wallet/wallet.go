package wallet

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const AddressPrefix = "sx1"

type Wallet struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func NewWallet() (*Wallet, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Wallet{PrivateKey: privateKey, PublicKey: publicKey}, nil
}

func (w *Wallet) Address() string {
	sum := sha256.Sum256(w.PublicKey)
	return AddressPrefix + hex.EncodeToString(sum[:20])
}

func AddressFromPublicKey(publicKey ed25519.PublicKey) (string, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return "", errors.New("invalid public key length")
	}
	sum := sha256.Sum256(publicKey)
	return AddressPrefix + hex.EncodeToString(sum[:20]), nil
}
