package main

import (
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/executor"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/desmos-labs/juno/version"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/auth"
	"github.com/forbole/bdjuno/x/bank"
	"github.com/forbole/bdjuno/x/staking"
	"github.com/forbole/bdjuno/x/supply"
	"github.com/go-co-op/gocron"
)

func main() {
	// Register all the modules to be handled
	SetupModules()

	// Build the executor
	prefix := "desmos" // TODO: Get this from a command
	rootCmd := executor.BuildRootCmd("bdjuno", SetupConfig(prefix))
	rootCmd.AddCommand(
		version.GetVersionCmd(),
		parse.GetParseCmd(simapp.MakeCodec(), database.Builder),
	)

	command := config.PrepareMainCmd(rootCmd)

	// Run the commands and panic on any error
	err := command.Execute()
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

func SetupModules() {
	// Register genesis handlers
	worker.RegisterGenesisHandler(auth.GenesisHandler)
	worker.RegisterGenesisHandler(staking.GenesisHandler)

	// Register block handlers
	worker.RegisterBlockHandler(staking.BlockHandler)

	// Register msg handlers
	worker.RegisterMsgHandler(staking.MsgHandler)
	worker.RegisterMsgHandler(bank.MsgHandler)
	worker.RegisterMsgHandler(staking.MsgHandler)

	// Register periodic operations
	scheduler := gocron.NewScheduler(time.UTC)
	parse.RegisterAdditionalOperation(staking.PeriodicStakingOperations(scheduler))
	parse.RegisterAdditionalOperation(auth.PeriodicAuthOperations(scheduler))
	parse.RegisterAdditionalOperation(supply.PeriodicSupplyOperations(scheduler))

	scheduler.StartAsync()
}
