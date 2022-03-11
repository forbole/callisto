package utils

import (
	"fmt"

	"github.com/forbole/juno/v3/node"
	"github.com/forbole/juno/v3/types/config"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ReadGenesis reads the genesis data based on the given config
func ReadGenesis(config config.Config, node node.Node) (*tmtypes.GenesisDoc, error) {
	if config.Parser.GenesisFilePath != "" {
		return readGenesisFromFilePath(config.Parser.GenesisFilePath)
	}

	return readGenesisFromNode(node)
}

func readGenesisFromFilePath(path string) (*tmtypes.GenesisDoc, error) {
	bz, err := tmos.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis file: %s", err)
	}

	var genDoc tmtypes.GenesisDoc
	err = tmjson.Unmarshal(bz, &genDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis doc: %s", err)
	}

	return &genDoc, nil
}

func readGenesisFromNode(node node.Node) (*tmtypes.GenesisDoc, error) {
	response, err := node.Genesis()
	if err != nil {
		return nil, fmt.Errorf("failed to get genesis: %s", err)
	}

	return response.Genesis, nil
}
