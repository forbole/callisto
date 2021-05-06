package bigdipper

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/db"
	jmodules "github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/desmos-labs/juno/modules/registrar"
	juno "github.com/desmos-labs/juno/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	"github.com/forbole/bdjuno/modules/bigdipper/bank"
	"github.com/forbole/bdjuno/modules/bigdipper/consensus"
	"github.com/forbole/bdjuno/modules/bigdipper/distribution"
	"github.com/forbole/bdjuno/modules/bigdipper/gov"
	"github.com/forbole/bdjuno/modules/bigdipper/mint"
	"github.com/forbole/bdjuno/modules/bigdipper/modules"
	"github.com/forbole/bdjuno/modules/bigdipper/slashing"
	"github.com/forbole/bdjuno/modules/bigdipper/staking"
	"github.com/forbole/bdjuno/modules/common/auth"
	"github.com/forbole/bdjuno/modules/common/pricefeed"
	"github.com/forbole/bdjuno/modules/common/utils"
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
	bigDipperBd := bigdipperdb.Cast(db)
	return []jmodules.Module{
		messages.NewModule(parser, encodingConfig.Marshaler, db),
		auth.NewModule(parser, encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		bank.NewModule(parser, encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		consensus.NewModule(cp, bigDipperBd),
		distribution.NewModule(utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		gov.NewModule(encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		mint.NewModule(utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		modules.NewModule(cfg, bigDipperBd),
		pricefeed.NewModule(encodingConfig, bigDipperBd),
		slashing.NewModule(utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		staking.NewModule(encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
	}
}
