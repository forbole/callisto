package modules

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/db"
	jmodules "github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/desmos-labs/juno/modules/registrar"
	juno "github.com/desmos-labs/juno/types"

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
	return []jmodules.Module{
		messages.NewModule(parser, encodingConfig.Marshaler, db),
		auth.NewModule(parser, encodingConfig, client.MustCreateGrpcConnection(cfg), bigDipperBd),
		bank.NewModule(parser, encodingConfig, client.MustCreateGrpcConnection(cfg), bigDipperBd),
		consensus.NewModule(cp, bigDipperBd),
		distribution.NewModule(client.MustCreateGrpcConnection(cfg), bigDipperBd),
		gov.NewModule(encodingConfig, client.MustCreateGrpcConnection(cfg), bigDipperBd),
		mint.NewModule(client.MustCreateGrpcConnection(cfg), bigDipperBd),
		modules.NewModule(cfg, bigDipperBd),
		pricefeed.NewModule(encodingConfig, bigDipperBd),
		slashing.NewModule(client.MustCreateGrpcConnection(cfg), bigDipperBd),
		staking.NewModule(encodingConfig, client.MustCreateGrpcConnection(cfg), bigDipperBd),
	}
}
