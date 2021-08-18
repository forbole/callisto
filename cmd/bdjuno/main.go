package main

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/desmos-labs/juno/cmd"
	parsecmd "github.com/desmos-labs/juno/cmd/parse"
	"github.com/desmos-labs/juno/modules/messages"

	"github.com/forbole/bdjuno/types/config"

	likeapp "github.com/likecoin/likechain/app"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules"

	likeapp "github.com/likecoin/likechain/app"
)

func main() {
	parseCfg := parsecmd.NewConfig().
		WithDBBuilder(database.Builder).
		WithConfigParser(config.Parser).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	cfg := cmd.NewConfig("bdjuno").
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
		likeapp.ModuleBasics,
	}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		messages.CosmosMessageAddressesParser,
	)
}
