package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/modules/registrar"
	"github.com/forbole/bdjuno/x/mint"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		pricefeed.Module{},
		staking.Module{},
		supply.Module{},
	)

	// Build the executor
	prefix := "desmos" // TODO: Get this from a command
	executor := cmd.BuildDefaultExecutor("bdjuno", SetupConfig(prefix), MakeCodec, database.Builder)

	// Run the commands and panic on any error
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

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

func MakeCodec() *codec.Codec {
	cdc := simapp.MakeCodec()
	return cdc
}
