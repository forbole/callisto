package main

import (
	"github.com/forbole/juno/v6/cmd"
	initcmd "github.com/forbole/juno/v6/cmd/init"
	parsetypes "github.com/forbole/juno/v6/cmd/parse/types"
	startcmd "github.com/forbole/juno/v6/cmd/start"
	"github.com/forbole/juno/v6/modules/messages"

	migratecmd "github.com/forbole/callisto/v4/cmd/migrate"
	parsecmd "github.com/forbole/callisto/v4/cmd/parse"
	"github.com/forbole/callisto/v4/utils"

	"github.com/forbole/callisto/v4/types/config"

	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules"
)

func main() {
	initCfg := initcmd.NewConfig().
		WithConfigCreator(config.Creator)

	cdc := utils.GetCodec()
	parseCfg := parsetypes.NewConfig().
		WithDBBuilder(database.Builder(cdc)).
		WithRegistrar(modules.NewRegistrar(getAddressesParser(), cdc))

	cfg := cmd.NewConfig("callisto").
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

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		messages.CosmosMessageAddressesParser,
	)
}
