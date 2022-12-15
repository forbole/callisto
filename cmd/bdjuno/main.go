package main

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	gaiaapp "github.com/cosmos/gaia/v7/app"
	desmosapp "github.com/desmos-labs/desmos/v4/app"
	"github.com/forbole/juno/v4/cmd"
	initcmd "github.com/forbole/juno/v4/cmd/init"
	migratecmd "github.com/forbole/juno/v4/cmd/migrate"
	parsetypes "github.com/forbole/juno/v4/cmd/parse/types"
	startcmd "github.com/forbole/juno/v4/cmd/start"
	"github.com/forbole/juno/v4/modules/messages"

	parsecmd "github.com/forbole/bdjuno/v3/cmd/parse"
	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules"
	"github.com/forbole/bdjuno/v3/types/config"
)

func main() {
	initCfg := initcmd.NewConfig().
		WithConfigCreator(config.Creator)

	parseCfg := parsetypes.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	cfg := cmd.NewConfig("bdjuno").
		WithInitConfig(initCfg).
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := cmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		cmd.VersionCmd(),
		initcmd.NewInitCmd(cfg.GetInitConfig()),
		parsecmd.NewParseCmd(cfg.GetParseConfig()),
		migratecmd.NewMigrateCmd(cfg.GetName(), cfg.GetParseConfig()),
		startcmd.NewStartCmd(cfg.GetParseConfig()),
	)

	executor := cmd.PrepareRootCmd(cfg.GetName(), rootCmd)
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
		gaiaapp.ModuleBasics,
		desmosapp.ModuleBasics,
	}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		desmosMessageAddressesParser,
		messages.CosmosMessageAddressesParser,
	)
}
