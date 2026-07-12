package config

import (
	"context"
	"os"
)

type Config struct {
	RPCAddress string
	P2PAddress string
	DataDir    string
	NetworkID  string
}

func New() *Config {
	return &Config{
		RPCAddress: getEnv("SAYANOX_RPC", ":8545"),
		P2PAddress: getEnv("SAYANOX_P2P", ":30303"),
		DataDir:    getEnv("SAYANOX_DATA", "./data"),
		NetworkID:  getEnv("SAYANOX_NETWORK", "sayanox-mainnet"),
	}
}

func (c *Config) Name() string {
	return "Config"
}

func (c *Config) Start(ctx context.Context) error {
	return nil
}

func (c *Config) Stop(ctx context.Context) error {
	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
