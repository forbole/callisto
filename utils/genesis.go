package utils

import (
	"fmt"

	cbftjson "github.com/cometbft/cometbft/libs/json"
	cbftos "github.com/cometbft/cometbft/libs/os"
	cbfttypes "github.com/cometbft/cometbft/types"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/types/config"
)

// ReadGenesis reads the genesis data based on the given config
func ReadGenesis(config config.Config, node node.Node) (*cbfttypes.GenesisDoc, error) {
	if config.Parser.GenesisFilePath != "" {
		return readGenesisFromFilePath(config.Parser.GenesisFilePath)
	}

	return readGenesisFromNode(node)
}

func readGenesisFromFilePath(path string) (*cbfttypes.GenesisDoc, error) {
	bz, err := cbftos.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis file: %s", err)
	}

	var genDoc cbfttypes.GenesisDoc
	err = cbftjson.Unmarshal(bz, &genDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis doc: %s", err)
	}

	return &genDoc, nil
}

func readGenesisFromNode(node node.Node) (*cbfttypes.GenesisDoc, error) {
	response, err := node.Genesis()
	if err != nil {
		return nil, fmt.Errorf("failed to get genesis: %s", err)
	}

	return response.Genesis, nil
}
