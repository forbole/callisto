package main

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/desmos-labs/juno/config"
	dbbuilder "github.com/desmos-labs/juno/db/builder"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse"
	djuno "github.com/desmos-labs/juno/types"
	"github.com/desmos-labs/juno/version"
)

func main() {
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
