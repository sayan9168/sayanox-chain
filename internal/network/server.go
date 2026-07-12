package network

import (
	"fmt"
	"net"
)

type Server struct {
	Host string
	Port int
}

func NewServer(host string, port int) *Server {
	return &Server{
		Host: host,
		Port: port,
	}
}

func (s *Server) Start() error {

	address := fmt.Sprintf(
		"%s:%d",
		s.Host,
		s.Port,
	)

	listener, err := net.Listen(
		"tcp",
		address,
	)

	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Println(
		"Sayanox node listening on",
		address,
	)

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	fmt.Println(
		"New peer connected:",
		conn.RemoteAddr(),
	)
}
