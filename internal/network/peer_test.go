package network

import (
	"net"
	"testing"
)

func TestPeerManagerReplacesDuplicate(t *testing.T) {
	pm := NewPeerManager()
	left, _ := net.Pipe()
	right, _ := net.Pipe()
	defer left.Close()
	defer right.Close()

	pm.AddPeer(&Peer{ID: "peer1", Conn: left})
	pm.AddPeer(&Peer{ID: "peer1", Conn: right})

	peers := pm.GetPeers()
	if len(peers) != 1 || peers[0].Conn != right {
		t.Fatalf("peer replacement failed")
	}
}
