package main

import (
	"github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/modules/registrar"

	"github.com/forbole/bdjuno/x/mint"
	"github.com/forbole/bdjuno/x/modules"

	"github.com/cosmos/cosmos-sdk/simapp"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/auth"
	"github.com/forbole/bdjuno/x/bank"
	"github.com/forbole/bdjuno/x/consensus"
	"github.com/forbole/bdjuno/x/distribution"
	"github.com/forbole/bdjuno/x/gov"
	"github.com/forbole/bdjuno/x/pricefeed"
	"github.com/forbole/bdjuno/x/staking"
	"github.com/forbole/bdjuno/x/supply"
)

func main() {
	// Register all the modules to be handled
	registrar.RegisterModules(
		auth.Module{},
		bank.Module{},
		consensus.Module{},
		distribution.Module{},
		gov.Module{},
		mint.Module{},
		modules.Module{},
		pricefeed.Module{},
		staking.Module{},
		supply.Module{},
	)

	executor := cmd.BuildDefaultExecutor("bdjuno", config.DefaultSetup, simapp.MakeCodec, database.Builder)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
