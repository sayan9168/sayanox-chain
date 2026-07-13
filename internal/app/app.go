package app

import (
	"context"

	"github.com/sayan9168/sayanox-chain/internal/node"
)

type App struct {
	node *node.Node
}

func New() *App {
	return &App{
		node: node.New(),
	}
}

func (a *App) Start(ctx context.Context) error {
	return a.node.Start(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	return a.node.Stop(ctx)
}

func (a *App) Node() *node.Node {
	return a.node
}
