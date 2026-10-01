package app

import (
	"context"
	"fmt"

	"github.com/sayan9168/sayanox-chain/configs"
	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/mempool"
	"github.com/sayan9168/sayanox-chain/internal/network"
	"github.com/sayan9168/sayanox-chain/internal/node"
	"github.com/sayan9168/sayanox-chain/internal/rpc"
	"github.com/sayan9168/sayanox-chain/internal/storage"
)

type App struct {
	config *configs.Config
	node   *node.Node
	store  *storage.ChainStore
}

func New() (*App, error) {
	cfg := configs.DefaultConfig()

	db, err := storage.Open(cfg.DataDir)
	if err != nil {
		return nil, err
	}

	store := storage.NewChainStore(db)
	chain, err := blockchain.NewPersistentChain(store)
	if err != nil {
		_ = store.Close()
		return nil, err
	}

	pool := mempool.NewMempool()
	networkServer := network.NewServer("0.0.0.0", cfg.P2PPort)
	rpcServer := rpc.NewServer(chain, pool, fmt.Sprintf(":%d", cfg.RPCPort))

	n := node.New()
	n.Register(networkServer)
	n.Register(rpcServer)

	return &App{
		config: cfg,
		node:   n,
		store:  store,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	return a.node.Start(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	err := a.node.Stop(ctx)
	if closeErr := a.store.Close(); err == nil {
		err = closeErr
	}
	return err
}

func (a *App) Config() *configs.Config {
	return a.config
}
