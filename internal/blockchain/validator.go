package blockchain

import "errors"

func ValidateBlock(
	current *Block,
	previous *Block,
) error {

	if current == nil {
		return errors.New("block is nil")
	}

	if previous != nil {
		if current.PreviousHash != previous.Hash {
			return errors.New("invalid previous hash")
		}

		if current.Height != previous.Height+1 {
			return errors.New("invalid block height")
		}
	}

	calculatedHash := current.CalculateHash()

	if current.Hash != calculatedHash {
		return errors.New("invalid block hash")
	}

	return nil
}
