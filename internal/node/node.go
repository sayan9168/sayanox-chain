package node

import (
	"context"

	"github.com/sayan9168/sayanox-chain/internal/service"
)

type Node struct {
	services []service.Service
}

func New() *Node {
	return &Node{
		services: make([]service.Service, 0),
	}
}

func (n *Node) Register(s service.Service) {
	n.services = append(n.services, s)
}

func (n *Node) Start(ctx context.Context) error {
	for _, s := range n.services {
		if err := s.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) Stop(ctx context.Context) error {
	for i := len(n.services) - 1; i >= 0; i-- {
		if err := n.services[i].Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
