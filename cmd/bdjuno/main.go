package main

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/config"
	dbbuilder "github.com/desmos-labs/juno/db/builder"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/version"
	modules "github.com/forbole/bdjuno/x"
	"github.com/forbole/bdjuno/x/staking"
)

func SetupConfig(prefix string) func(cfg *sdk.Config) {
	return func(cfg *sdk.Config) {
		cfg.SetBech32PrefixForAccount(
			prefix,
			prefix+sdk.PrefixPublic,
		)
		cfg.SetBech32PrefixForValidator(
			prefix+sdk.PrefixValidator+sdk.PrefixOperator,
			prefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
		)
		cfg.SetBech32PrefixForConsensusNode(
			prefix+sdk.PrefixValidator+sdk.PrefixConsensus,
			prefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
		)
	}
}

func main() {

	// Register all the modules to be handled
	modules.RegisterModule(modules.NewModule("staking", staking.Fetcher, staking.Handler))

	// Register the block handler that allows all modules to be properly updated
	worker.RegisterBlockHandler(modules.HandleModules)

	// Build the executor
	prefix := "desmos" // TODO: Get this from a command
	rootCmd := executor.BuildRootCmd("bdjuno", SetupConfig(prefix))
	rootCmd.AddCommand(
		version.GetVersionCmd(),
		parse.GetParseCmd(simapp.MakeCodec(), dbbuilder.DatabaseBuilder),
	)

	command := config.PrepareMainCmd(rootCmd)

	// Run the commands and panic on any error
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
