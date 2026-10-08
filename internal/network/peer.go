package network

import (
	"net"
	"sync"
	"time"
)

type Peer struct {
	ID       string
	Address  string
	Version  int
	Conn     net.Conn
	LastSeen time.Time
}

type PeerManager struct {
	Peers map[string]*Peer
	mu    sync.RWMutex
}

func NewPeerManager() *PeerManager {
	return &PeerManager{Peers: make(map[string]*Peer)}
}

func (pm *PeerManager) AddPeer(peer *Peer) {
	if peer == nil || peer.ID == "" { return }
	pm.mu.Lock()
	defer pm.mu.Unlock()
	peer.LastSeen = time.Now()
	if old, ok := pm.Peers[peer.ID]; ok && old.Conn != nil && old.Conn != peer.Conn {
		_ = old.Conn.Close()
	}
	pm.Peers[peer.ID] = peer
}

func (pm *PeerManager) RemovePeer(id string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if peer, ok := pm.Peers[id]; ok && peer.Conn != nil { _ = peer.Conn.Close() }
	delete(pm.Peers, id)
}

func (pm *PeerManager) GetPeers() []*Peer {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	list := make([]*Peer, 0, len(pm.Peers))
	for _, peer := range pm.Peers { list = append(list, peer) }
	return list
}
