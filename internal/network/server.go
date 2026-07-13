package network

import (
	"context"
	"fmt"
	"net"
)

type Server struct {
	Host     string
	Port     int
	listener net.Listener
}

func NewServer(host string, port int) *Server {
	return &Server{
		Host: host,
		Port: port,
	}
}

func (s *Server) Name() string {
	return "network"
}

func (s *Server) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

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

		go handleConnection(conn)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New peer connected:", conn.RemoteAddr())

	// TODO:
	// - Handshake
	// - Version exchange
	// - Peer authentication
	// - Message decoding
	// - Block sync
	// - Transaction sync
}
