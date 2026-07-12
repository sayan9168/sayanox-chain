package network

import (
	"sync"
	"time"
)

type Peer struct {
	ID        string
	Address   string
	LastSeen  time.Time
}

type PeerManager struct {
	Peers map[string]*Peer
	mu    sync.RWMutex
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		Peers: make(map[string]*Peer),
	}
}

func (pm *PeerManager) AddPeer(peer *Peer) {

	pm.mu.Lock()
	defer pm.mu.Unlock()

	peer.LastSeen = time.Now()

	pm.Peers[peer.ID] = peer
}

func (pm *PeerManager) RemovePeer(id string) {

	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.Peers, id)
}

func (pm *PeerManager) GetPeers() []*Peer {

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	list := make([]*Peer, 0)

	for _, peer := range pm.Peers {
		list = append(list, peer)
	}

	return list
}
