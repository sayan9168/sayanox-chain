package blockchain

import "errors"

func ValidateBlock(current *Block, previous *Block) error {
	if current == nil { return errors.New("block is nil") }
	if previous != nil {
		if current.PreviousHash != previous.Hash { return errors.New("invalid previous hash") }
		if current.Height != previous.Height+1 { return errors.New("invalid block height") }
	}
	if current.Hash != current.CalculateHash() { return errors.New("invalid block hash") }
	for _, tx := range current.Transactions {
		if tx == nil { return errors.New("nil transaction") }
	}
	return nil
}
