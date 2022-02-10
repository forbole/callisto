package main

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/forbole/juno/v2/cmd"
	initcmd "github.com/forbole/juno/v2/cmd/init"
	parsecmd "github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/modules/messages"

	actionscmd "github.com/forbole/bdjuno/v2/cmd/actions"
	fixcmd "github.com/forbole/bdjuno/v2/cmd/fix"
	migratecmd "github.com/forbole/bdjuno/v2/cmd/migrate"
	parsegenesiscmd "github.com/forbole/bdjuno/v2/cmd/parse-genesis"

	"github.com/forbole/bdjuno/v2/types/config"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"

	gaiaapp "github.com/cosmos/gaia/v6/app"
	desmosapp "github.com/desmos-labs/desmos/v2/app"
)

func main() {
	parseCfg := parsecmd.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	cfg := cmd.NewConfig("bdjuno").
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := cmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		cmd.VersionCmd(),
		initcmd.InitCmd(cfg.GetInitConfig()),
		parsecmd.ParseCmd(cfg.GetParseConfig()),
		migratecmd.NewMigrateCmd(),
		fixcmd.NewFixCmd(cfg.GetParseConfig()),
		parsegenesiscmd.NewParseGenesisCmd(cfg.GetParseConfig()),
		actionscmd.NewActionsCmd(cfg.GetParseConfig()),
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
