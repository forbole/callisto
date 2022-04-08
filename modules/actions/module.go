package actions

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/node"
	nodeconfig "github.com/forbole/juno/v3/node/config"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/forbole/juno/v3/types/config"

	modulestypes "github.com/forbole/bdjuno/v2/modules/types"
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

func NewModule(cfg config.Config, encodingConfig *params.EncodingConfig) *Module {
	actionsCfg, err := ParseConfig(cfg.GetBytes())
	if err != nil {
		panic(err)
	}

	// Build the node
	junoNode, err := remote.NewNode(actionsCfg.Node, encodingConfig.Marshaler)
	if err != nil {
		panic(err)
	}

	// Build the sources
	sources, err := modulestypes.BuildSources(nodeconfig.NewConfig("remote", actionsCfg.Node), encodingConfig)
	if err != nil {
		panic(err)
	}

	return &Module{
		cfg:     actionsCfg,
		node:    junoNode,
		sources: sources,
	}
}

func (m *Module) Name() string {
	return "actions"
}
