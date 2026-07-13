package app

import (
	"context"
	"fmt"

	"github.com/sayan9168/sayanox-chain/configs"
	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/network"
	"github.com/sayan9168/sayanox-chain/internal/node"
	"github.com/sayan9168/sayanox-chain/internal/rpc"
)

type App struct {
	config *configs.Config
	node   *node.Node
}

func New() *App {
	cfg := configs.DefaultConfig()

	chain := blockchain.NewChain()

	networkServer := network.NewServer(
		"0.0.0.0",
		cfg.P2PPort,
	)

	rpcServer := rpc.NewServer(
		chain,
		fmt.Sprintf(":%d", cfg.RPCPort),
	)

	n := node.New()

	n.Register(networkServer)
	n.Register(rpcServer)

	return &App{
		config: cfg,
		node:   n,
	}
}

func (a *App) Start(ctx context.Context) error {
	return a.node.Start(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	return a.node.Stop(ctx)
}

func (a *App) Config() *configs.Config {
	return a.config
}
