package network

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"net"
	"sync"

	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

type Server struct {
	Host       string
	Port       int
	ChainID    string
	listener   net.Listener
	privateKey ed25519.PrivateKey
	peers      *PeerManager
	mu         sync.RWMutex
}

func NewServer(host string, port int) *Server {
	w, err := wallet.NewWallet()
	if err != nil {
		panic(err)
	}
	return &Server{
		Host: host, Port: port, ChainID: "sayanox-mainnet",
		privateKey: w.PrivateKey, peers: NewPeerManager(),
	}
}

func (s *Server) Name() string { return "network" }

func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil { return err }
	s.listener = listener
	fmt.Println("Sayanox node listening on", address)

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) Dial(ctx context.Context, address string) error {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil { return err }
	info, err := performHandshake(conn, s.ChainID, s.privateKey, true)
	if err != nil {
		_ = conn.Close()
		return err
	}
	s.peers.AddPeer(&Peer{
		ID: info.PeerID, Address: address, Version: info.Version, Conn: conn,
	})
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.listener != nil { return s.listener.Close() }
	for _, peer := range s.peers.GetPeers() {
		s.peers.RemovePeer(peer.ID)
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	info, err := performHandshake(conn, s.ChainID, s.privateKey, false)
	if err != nil {
		_ = conn.Close()
		return
	}
	s.peers.AddPeer(&Peer{
		ID: info.PeerID, Address: conn.RemoteAddr().String(),
		Version: info.Version, Conn: conn,
	})
}
