package main

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/modules/messages"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x"
)

func main() {
	// Build the executor
	executor := cmd.BuildDefaultExecutor(
		"bdjuno",
		x.NewModulesRegistrar(
			messages.CosmosMessageAddressesParser,
		),
		juno.DefaultSetup,
		simapp.MakeTestEncodingConfig,
		database.Builder,
	)

	// Run the command
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
