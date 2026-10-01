package transaction

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

type Transaction struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    *Amount `json:"amount"`
	Fee       *Amount `json:"fee"`
	Timestamp int64   `json:"timestamp"`
	Nonce     uint64  `json:"nonce"`
	Hash      string  `json:"hash"`
	Signature string  `json:"signature"`
	PublicKey string  `json:"public_key"`
}

type signingPayload struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    string `json:"amount"`
	Fee       string `json:"fee"`
	Timestamp int64  `json:"timestamp"`
	Nonce     uint64 `json:"nonce"`
	PublicKey string `json:"public_key"`
}

func NewTransaction(from, to string, amount, fee *Amount) *Transaction {
	tx := &Transaction{
		From: from, To: to, Amount: amount, Fee: fee,
		Timestamp: time.Now().Unix(),
	}
	tx.Hash = tx.CalculateHash()
	return tx
}

func NewSignedTransaction(w *wallet.Wallet, to string, amount, fee *Amount, nonce uint64) (*Transaction, error) {
	if w == nil {
		return nil, errors.New("nil wallet")
	}
	from := w.Address()
	tx := &Transaction{
		From: from, To: to, Amount: amount, Fee: fee,
		Timestamp: time.Now().Unix(), Nonce: nonce,
		PublicKey: hex.EncodeToString(w.PublicKey),
	}
	tx.Hash = tx.CalculateHash()
	tx.Sign(w.PrivateKey)
	return tx, nil
}

func (tx *Transaction) signingBytes() ([]byte, error) {
	if tx == nil || tx.Amount == nil || tx.Fee == nil {
		return nil, errors.New("invalid transaction")
	}
	p := signingPayload{
		From: tx.From, To: tx.To, Amount: tx.Amount.String(),
		Fee: tx.Fee.String(), Timestamp: tx.Timestamp, Nonce: tx.Nonce,
		PublicKey: tx.PublicKey,
	}
	return json.Marshal(p)
}

func (tx *Transaction) CalculateHash() string {
	data, err := tx.signingBytes()
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (tx *Transaction) Sign(privateKey ed25519.PrivateKey) {
	tx.Hash = tx.CalculateHash()
	tx.Signature = hex.EncodeToString(ed25519.Sign(privateKey, []byte(tx.Hash)))
}

func (tx *Transaction) Verify() bool {
	if tx == nil || tx.Hash == "" || tx.Signature == "" || tx.PublicKey == "" {
		return false
	}
	publicKey, err := hex.DecodeString(tx.PublicKey)
	if err != nil || len(publicKey) != ed25519.PublicKeySize {
		return false
	}
	signature, err := hex.DecodeString(tx.Signature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return false
	}
	if tx.CalculateHash() != tx.Hash {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(publicKey), []byte(tx.Hash), signature)
}

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

func publicKeyBytes(encoded string) ed25519.PublicKey {
	b, err := hex.DecodeString(encoded)
	if err != nil {
		return nil
	}
	return ed25519.PublicKey(b)
}
