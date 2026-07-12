package wallet

import (
	"encoding/hex"
	"encoding/json"
	"os"

	"crypto/ed25519"
)

type KeyStore struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

func SaveWallet(
	path string,
	w *Wallet,
	address string,
) error {

	store := KeyStore{
		PrivateKey: hex.EncodeToString(w.PrivateKey),
		PublicKey:  hex.EncodeToString(w.PublicKey),
		Address:    address,
	}

	data, err := json.MarshalIndent(
		store,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0600,
	)
}

func LoadWallet(path string) (*Wallet, string, error) {

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, "", err
	}

	var store KeyStore

	err = json.Unmarshal(
		data,
		&store,
	)

	if err != nil {
		return nil, "", err
	}

	privateKeyBytes, err := hex.DecodeString(
		store.PrivateKey,
	)

	if err != nil {
		return nil, "", err
	}

	publicKeyBytes, err := hex.DecodeString(
		store.PublicKey,
	)

	if err != nil {
		return nil, "", err
	}

	return &Wallet{
		PrivateKey: ed25519.PrivateKey(privateKeyBytes),
		PublicKey:  ed25519.PublicKey(publicKeyBytes),
	}, store.Address, nil
}
