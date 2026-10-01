package main

import (
	"fmt"
	"log"

	"github.com/sayan9168/sayanox-chain/configs"
	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/network"
	"github.com/sayan9168/sayanox-chain/internal/rpc"
)

func main() {
	config := configs.DefaultConfig()
	fmt.Println("Starting", config.NetworkName, config.Version)

	chain := blockchain.NewChain()
	networkServer := network.NewServer("0.0.0.0", config.P2PPort)
	rpcServer := rpc.NewServer(chain, fmt.Sprintf(":%d", config.RPCPort))

	go func() {
		if err := networkServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Printf("RPC server listening on %s\n", rpcServer.Port)
	if err := rpcServer.Start(); err != nil {
		log.Fatal(err)
	}
}
