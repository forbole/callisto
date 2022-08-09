package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
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

	sdk "github.com/cosmos/cosmos-sdk/types"
	gammtypes "github.com/osmosis-labs/osmosis/v10/x/gamm/pool-models/balancer"

	gaiaapp "github.com/cosmos/gaia/v7/app"
	osmosisapp "github.com/osmosis-labs/osmosis/v10/app"
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
		osmosisapp.ModuleBasics,
	}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		OsmoMessageAddressesParser,
		messages.CosmosMessageAddressesParser,
	)
}

// OsmoMessageAddressesParser returns the list of all the accounts involved in the given
// message if it's related to the x/gamm module
func OsmoMessageAddressesParser(_ codec.Codec, cosmosMsg sdk.Msg) ([]string, error) {

	if msg, ok := cosmosMsg.(*gammtypes.MsgCreateBalancerPool); ok {
		return []string{msg.Sender, msg.FuturePoolGovernor}, nil
	}

	return nil, messages.MessageNotSupported(cosmosMsg)
}
