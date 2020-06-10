package main

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/desmos-labs/juno/config"
	dbbuilder "github.com/desmos-labs/juno/db/builder"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/worker"
	djuno "github.com/desmos-labs/juno/types"
	"github.com/desmos-labs/juno/version"
	modules "github.com/forbole/bdjuno/x"
	"github.com/forbole/bdjuno/x/staking"
)

func main() {
	// Register all the modules to be handled
	modules.RegisterModule(modules.NewModule("staking", staking.Fetcher, staking.Handler))

	// Register the block handler that allows all modules to be properly updated
	worker.RegisterBlockHandler(modules.HandleModules)

	// Build the executor
	rootCmd := executor.BuildRootCmd("bdjuno", djuno.EmptySetup)
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
