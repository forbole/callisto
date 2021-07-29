package main

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	desmosapp "github.com/desmos-labs/desmos/app"
	"github.com/desmos-labs/juno/cmd"
	initcmd "github.com/desmos-labs/juno/cmd/init"
	parsecmd "github.com/desmos-labs/juno/cmd/parse"

	"github.com/forbole/bdjuno/types/config"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules"
)

func main() {
	// Setup the config
	initCfg := initcmd.NewConfig().
		WithConfigFlagSetup(config.SetupConfigFlags).
		WithConfigCreator(config.CreateConfig)

	parseCfg := parsecmd.NewConfig().
		WithConfigParser(config.ParseConfig).
		WithRegistrar(modules.NewRegistrar()).
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers()))

	cfg := cmd.NewConfig("bdjuno").
		WithInitConfig(initCfg).
		WithParseConfig(parseCfg)

	// Run the command
	executor := cmd.BuildDefaultExecutor(cfg)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

// getBasicManagers returns the various basic managers that are used to
// register the encoding to support custom messages
func getBasicManagers() []module.BasicManager {
	return []module.BasicManager{
		simapp.ModuleBasics,
		desmosapp.ModuleBasics,
	}
}
