package genesis

import (
	_ "embed"
	"encoding/json"
)

//go:embed genesis.json
var embeddedGenesis []byte

func Embedded() (*Genesis, error) {
	var g Genesis
	if err := json.Unmarshal(embeddedGenesis, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func MustEmbedded() *Genesis {
	g, err := Embedded()
	if err != nil {
		panic(err)
	}
	return g
}
