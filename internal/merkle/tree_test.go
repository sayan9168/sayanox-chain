package merkle

import "testing"

func TestMerkleRoot(t *testing.T) {
	root := Root([]string{
		"tx1",
		"tx2",
		"tx3",
		"tx4",
	})

	if root == "" {
		t.Fatal("merkle root is empty")
	}
}
