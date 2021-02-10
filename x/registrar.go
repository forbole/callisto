package x

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/auth"
	"github.com/forbole/bdjuno/x/bank"
	"github.com/forbole/bdjuno/x/consensus"
	"github.com/forbole/bdjuno/x/distribution"
	"github.com/forbole/bdjuno/x/gov"
	"github.com/forbole/bdjuno/x/mint"
	bmodules "github.com/forbole/bdjuno/x/modules"
	"github.com/forbole/bdjuno/x/pricefeed"
	"github.com/forbole/bdjuno/x/slashing"
	"github.com/forbole/bdjuno/x/staking"
	"github.com/forbole/bdjuno/x/utils"
)

// ModulesRegistrar represents the modules.Registrar that allows to register all custom BDJuno modules
type ModulesRegistrar struct {
	parser messages.MessageAddressesParser
}

// NewModulesRegistrar allows to build a new ModulesRegistrar instance
func NewModulesRegistrar(parser messages.MessageAddressesParser) *ModulesRegistrar {
	return &ModulesRegistrar{
		parser: parser,
	}
}

// BuildModules implements modules.Registrar
func (r *ModulesRegistrar) BuildModules(
	cfg *config.Config, encodingConfig *params.EncodingConfig, _ *sdk.Config, db db.Database, cp *client.Proxy,
) modules.Modules {
	bigDipperBd := database.Cast(db)
	return []modules.Module{
		messages.NewModule(r.parser, encodingConfig.Marshaler, db),
		auth.NewModule(encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		bank.NewModule(encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		consensus.NewModule(cp, bigDipperBd),
		distribution.NewModule(utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		gov.NewModule(encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		mint.NewModule(utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		bmodules.NewModule(cfg, bigDipperBd),
		pricefeed.NewModule(bigDipperBd),
		slashing.NewModule(utils.MustCreateGrpcConnection(cfg), bigDipperBd),
		staking.NewModule(encodingConfig, utils.MustCreateGrpcConnection(cfg), bigDipperBd),
	}
}
