package genesis

import (
	"encoding/json"
	"os"
)

type Genesis struct {
	ChainID       string   `json:"chain_id"`
	NetworkName   string   `json:"network_name"`
	Version       string   `json:"version"`
	GenesisTime   string   `json:"genesis_time"`
	InitialSupply uint64   `json:"initial_supply"`
	Symbol        string   `json:"symbol"`
	Decimals      uint8    `json:"decimals"`
	Validators    []string `json:"validators"`
}

func LoadGenesis(path string) (*Genesis, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var g Genesis
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}

	return &g, nil
}
