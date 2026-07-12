package node

import (
	"fmt"

	"github.com/sayan9168/sayanox-chain/internal/blockchain"
	"github.com/sayan9168/sayanox-chain/internal/consensus"
	"github.com/sayan9168/sayanox-chain/internal/network"
	"github.com/sayan9168/sayanox-chain/internal/rpc"
)

type Node struct {
	Chain     *blockchain.Chain
	Consensus *consensus.Engine
	Network   *network.Server
	RPC       *rpc.Server
}

func NewNode(
	chain *blockchain.Chain,
	consensus *consensus.Engine,
	network *network.Server,
	rpcServer *rpc.Server,
) *Node {

	return &Node{
		Chain:     chain,
		Consensus: consensus,
		Network:   network,
		RPC:       rpcServer,
	}
}

func (n *Node) Start() error {

	go func() {
		if err := n.Network.Start(); err != nil {
			fmt.Println("Network error:", err)
		}
	}()

	return n.RPC.Start()
}
