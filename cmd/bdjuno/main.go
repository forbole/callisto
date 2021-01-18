package main

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"

	"github.com/forbole/bdjuno/x/slashing"

	"github.com/forbole/bdjuno/x/mint"
	bmodules "github.com/forbole/bdjuno/x/modules"

	"github.com/cosmos/cosmos-sdk/simapp"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/auth"
	"github.com/forbole/bdjuno/x/bank"
	"github.com/forbole/bdjuno/x/consensus"
	"github.com/forbole/bdjuno/x/distribution"
	"github.com/forbole/bdjuno/x/gov"
	"github.com/forbole/bdjuno/x/pricefeed"
	"github.com/forbole/bdjuno/x/staking"
)

func main() {
	executor := cmd.BuildDefaultExecutor(
		"bdjuno",
		&ModulesRegistrar{},
		config.DefaultSetup,
		simapp.MakeTestEncodingConfig,
		database.Builder,
	)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

type ModulesRegistrar struct{}

func (*ModulesRegistrar) BuildModules(
	cfg *config.Config, encodingConfig *params.EncodingConfig, _ *sdk.Config, db db.Database, cp *client.Proxy,
) modules.Modules {
	bigDipperBd := database.Cast(db)
	return []modules.Module{
		auth.NewModule(encodingConfig, cp.GrpcConnection(), bigDipperBd),
		bank.NewModule(encodingConfig, cp.GrpcConnection(), bigDipperBd),
		consensus.NewModule(cp, bigDipperBd),
		distribution.NewModule(cp.GrpcConnection(), bigDipperBd),
		gov.NewModule(encodingConfig, cp.GrpcConnection(), bigDipperBd),
		mint.NewModule(cp.GrpcConnection(), bigDipperBd),
		bmodules.NewModule(cfg, bigDipperBd),
		pricefeed.NewModule(bigDipperBd),
		slashing.NewModule(cp.GrpcConnection(), bigDipperBd),
		staking.NewModule(encodingConfig, cp.GrpcConnection(), bigDipperBd),
	}
}
