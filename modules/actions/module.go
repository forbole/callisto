package actions

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/forbole/juno/v6/modules"
	"github.com/forbole/juno/v6/node"
	"github.com/forbole/juno/v6/node/builder"
	nodeconfig "github.com/forbole/juno/v6/node/config"
	"github.com/forbole/juno/v6/types/config"
	"github.com/rs/zerolog/log"

	modulestypes "github.com/forbole/callisto/v4/modules/types"
)

const (
	ModuleName = "actions"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

type Module struct {
	cfg     *Config
	node    node.Node
	sources *modulestypes.Sources
}

func NewModule(cfg config.Config, cdc codec.Codec, sources *modulestypes.Sources) *Module {
	bz, err := cfg.GetBytes()
	if err != nil {
		panic(err)
	}

	actionsCfg, err := ParseConfig(bz)
	if err != nil {
		panic(err)
	}

	nodeCfg := cfg.Node
	if actionsCfg.Node != nil {
		nodeCfg = nodeconfig.NewConfig(nodeconfig.TypeRemote, actionsCfg.Node)
	}

	var node node.Node
	if cfg.Node.Type == nodeconfig.TypeLocal {
		log.Err(fmt.Errorf("local node is not supported for actions module"))
	} else {
		// Build the node
		txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
		junoNode, err := builder.BuildNode(nodeCfg, txConfig, cdc)
		if err != nil {
			panic(err)
		}

		node = junoNode
	}

	return &Module{
		cfg:     actionsCfg,
		node:    node,
		sources: sources,
	}
}

func (m *Module) Name() string {
	return ModuleName
}
