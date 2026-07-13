package node

import (
	"context"

	"github.com/sayan9168/sayanox-chain/internal/service"
)

type Node struct {
	services []service.Service
	errors   chan error
}

func New() *Node {
	return &Node{
		services: make([]service.Service, 0),
		errors:   make(chan error, 16),
	}
}

func (n *Node) Register(s service.Service) {
	n.services = append(n.services, s)
}

func (n *Node) Start(ctx context.Context) error {
	for _, svc := range n.services {
		go func(s service.Service) {
			if err := s.Start(ctx); err != nil {
				select {
				case n.errors <- err:
				default:
				}
			}
		}(svc)
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

func (n *Node) Errors() <-chan error {
	return n.errors
}
