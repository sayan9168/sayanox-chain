package consensus

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

type Engine struct {
	PoS *PoS
}

func NewEngine() *Engine {
	return &Engine{PoS: NewPoS()}
}

func (e *Engine) Proposer(previous *blockchain.Block) (*Validator, error) {
	if e == nil || e.PoS == nil {
		return nil, errors.New("consensus engine not configured")
	}
	seed, err := ProposerSeed(previous)
	if err != nil {
		return nil, err
	}
	return e.PoS.SelectValidator(seed)
}

func (e *Engine) ValidateBlock(block *blockchain.Block, previous *blockchain.Block) error {
	if e == nil || e.PoS == nil {
		return errors.New("consensus engine not configured")
	}
	if block == nil || previous == nil {
		return errors.New("block and previous block are required")
	}
	if err := block.Validate(previous); err != nil {
		return err
	}
	if block.Proposer == "" || block.ProposerSignature == "" || block.ProposerPublicKey == "" {
		return errors.New("missing proposer authentication")
	}

	expected, err := e.Proposer(previous)
	if err != nil {
		return err
	}
	if block.Proposer != expected.Address {
		return errors.New("block proposer is not selected validator")
	}

	publicKey, err := hex.DecodeString(block.ProposerPublicKey)
	if err != nil || len(publicKey) != ed25519.PublicKeySize {
		return errors.New("invalid proposer public key")
	}
	if len(expected.PublicKey) != 0 && !bytes.Equal(expected.PublicKey, publicKey) {
		return errors.New("proposer public key mismatch")
	}

	address, err := wallet.AddressFromPublicKey(ed25519.PublicKey(publicKey))
	if err != nil || address != block.Proposer {
		return errors.New("proposer address mismatch")
	}

	signature, err := hex.DecodeString(block.ProposerSignature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return errors.New("invalid proposer signature")
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), []byte(block.SigningHash()), signature) {
		return errors.New("invalid proposer signature")
	}
	return nil
}

func ProposerSeed(previous *blockchain.Block) (*big.Int, error) {
	if previous == nil {
		return nil, errors.New("previous block is nil")
	}
	hash, err := hex.DecodeString(previous.Hash)
	if err != nil || len(hash) == 0 {
		return nil, errors.New("invalid previous block hash")
	}
	return new(big.Int).SetBytes(hash), nil
}
