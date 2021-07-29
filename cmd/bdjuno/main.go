package main

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	desmosapp "github.com/desmos-labs/desmos/app"
	"github.com/desmos-labs/juno/cmd"
	initcmd "github.com/desmos-labs/juno/cmd/init"
	parsecmd "github.com/desmos-labs/juno/cmd/parse"
	"github.com/desmos-labs/juno/modules/messages"

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
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

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

// getBasicManagers returns the various basic managers that are used to register the encoding to
// support custom messages.
// This should be edited by custom implementations if needed.
func getBasicManagers() []module.BasicManager {
	return []module.BasicManager{
		simapp.ModuleBasics,
		desmosapp.ModuleBasics,
	}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		messages.CosmosMessageAddressesParser,
		desmosMessageAddressesParser,
	)
}
