package merkle

import (
	"crypto/sha256"
	"encoding/hex"
)

func hash(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

func Root(leaves []string) string {
	if len(leaves) == 0 {
		return ""
	}

	nodes := make([]string, len(leaves))

	for i, leaf := range leaves {
		nodes[i] = hash(leaf)
	}

	for len(nodes) > 1 {
		var next []string

		for i := 0; i < len(nodes); i += 2 {

			left := nodes[i]
			right := left

			if i+1 < len(nodes) {
				right = nodes[i+1]
			}

			next = append(next, hash(left+right))
		}

		nodes = next
	}

	return nodes[0]
}
