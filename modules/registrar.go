package modules

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/db"
	jmodules "github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/desmos-labs/juno/modules/registrar"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/modules/history"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/auth"
	"github.com/forbole/bdjuno/modules/bank"
	"github.com/forbole/bdjuno/modules/consensus"
	"github.com/forbole/bdjuno/modules/distribution"
	"github.com/forbole/bdjuno/modules/gov"
	"github.com/forbole/bdjuno/modules/mint"
	"github.com/forbole/bdjuno/modules/modules"
	"github.com/forbole/bdjuno/modules/pricefeed"
	"github.com/forbole/bdjuno/modules/slashing"
	"github.com/forbole/bdjuno/modules/staking"
	"github.com/forbole/bdjuno/modules/utils"
)

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar() *Registrar {
	return &Registrar{}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(
	cfg juno.Config, encodingConfig *params.EncodingConfig, _ *sdk.Config, db db.Database, cp *client.Proxy,
) jmodules.Modules {
	parser := utils.AddressesParser
	bigDipperBd := database.Cast(db)
	grpcConnection := client.MustCreateGrpcConnection(cfg)

	authClient := authttypes.NewQueryClient(grpcConnection)
	bankClient := banktypes.NewQueryClient(grpcConnection)
	distrClient := distrtypes.NewQueryClient(grpcConnection)
	govClient := govtypes.NewQueryClient(grpcConnection)
	mintClient := minttypes.NewQueryClient(grpcConnection)
	slashingClient := slashingtypes.NewQueryClient(grpcConnection)
	stakingClient := stakingtypes.NewQueryClient(grpcConnection)

	return []jmodules.Module{
		messages.NewModule(parser, encodingConfig.Marshaler, db),
		auth.NewModule(parser, authClient, encodingConfig, bigDipperBd),
		bank.NewModule(parser, authClient, bankClient, encodingConfig, bigDipperBd),
		consensus.NewModule(cp, bigDipperBd),
		distribution.NewModule(distrClient, bigDipperBd),
		gov.NewModule(bankClient, govClient, stakingClient, encodingConfig, bigDipperBd),
		mint.NewModule(mintClient, bigDipperBd),
		modules.NewModule(cfg, bigDipperBd),
		pricefeed.NewModule(encodingConfig, bigDipperBd),
		slashing.NewModule(slashingClient, bigDipperBd),
		staking.NewModule(bankClient, stakingClient, encodingConfig, bigDipperBd),
		history.NewModule(parser, encodingConfig, bigDipperBd),
	}
}
