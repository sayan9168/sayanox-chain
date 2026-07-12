package configs

import "time"

type Config struct {
    NetworkName string
    ChainID     string
    Version     string

    BlockTime   time.Duration
    MaxBlockSize uint64

    RPCPort int
    P2PPort int

    DataDir string
}

func DefaultConfig() *Config {
    return &Config{
        NetworkName: "Sayanox Chain",
        ChainID:     "sayanox-mainnet",
        Version:     "v0.1.0-alpha",

        BlockTime:    2 * time.Second,
        MaxBlockSize: 2 * 1024 * 1024,

        RPCPort: 8545,
        P2PPort: 30303,

        DataDir: "./data",
    }
}
