package main

import (
	chainapp "git.ooo.ua/vipcoin/chain/app"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/forbole/juno/v3/cmd"
	initcmd "github.com/forbole/juno/v3/cmd/init"
	parsetypes "github.com/forbole/juno/v3/cmd/parse/types"
	startcmd "github.com/forbole/juno/v3/cmd/start"
	"github.com/forbole/juno/v3/modules/messages"

	migratecmd "github.com/forbole/bdjuno/v3/cmd/migrate"
	parsecmd "github.com/forbole/bdjuno/v3/cmd/parse"

	"github.com/forbole/bdjuno/v3/types/config"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules"

	gaiaapp "github.com/cosmos/gaia/v7/app"
	actionscmd "github.com/forbole/bdjuno/v2/cmd/actions"
	fixcmd "github.com/forbole/bdjuno/v2/cmd/fix"
	migratecmd "github.com/forbole/bdjuno/v2/cmd/migrate"
	parsegenesiscmd "github.com/forbole/bdjuno/v2/cmd/parse-genesis"
	vault "github.com/forbole/bdjuno/v2/config"
	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/types/config"
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
		vault.CheckVaultConfig(cfg.GetName(), parsecmd.ParseCmd(cfg.GetParseConfig())),
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
		chainapp.ModuleBasics,
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
