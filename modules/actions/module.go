package actions

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"github.com/forbole/juno/v4/node/builder"
	nodeconfig "github.com/forbole/juno/v4/node/config"
	"github.com/forbole/juno/v4/types/config"

	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
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

func NewModule(cfg config.Config, encodingConfig *params.EncodingConfig) *Module {
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

	// Build the node
	junoNode, err := builder.BuildNode(nodeCfg, encodingConfig)
	if err != nil {
		panic(err)
	}

	// Build the sources
	sources, err := modulestypes.BuildSources(nodeCfg, encodingConfig)
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
	return ModuleName
}
