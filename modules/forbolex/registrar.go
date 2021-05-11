package forbolex

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/registrar"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/modules/forbolex/distribution"

	forbolexdb "github.com/forbole/bdjuno/database/forbolex"
	"github.com/forbole/bdjuno/modules/common/auth"
	"github.com/forbole/bdjuno/modules/common/pricefeed"
	"github.com/forbole/bdjuno/modules/common/utils"
	"github.com/forbole/bdjuno/modules/forbolex/bank"
	"github.com/forbole/bdjuno/modules/forbolex/staking"
)

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents an implementation of registrar.Registrar that allows to register the
// modules compatible with ForboleX
type Registrar struct {
}

// NewRegistrar allows to build a new Registrar instance.
func NewRegistrar() *Registrar {
	return &Registrar{}
}

// BuildModules implements registrar.Registrar
func (r *Registrar) BuildModules(
	cfg juno.Config, encodingConfig *params.EncodingConfig, _ *sdk.Config, db db.Database, _ *client.Proxy,
) modules.Modules {
	parser := utils.AddressesParser
	forboleXDB := forbolexdb.Cast(db)
	return modules.Modules{
		auth.NewModule(parser, encodingConfig, utils.MustCreateGrpcConnection(cfg), forboleXDB),
		bank.NewModule(parser, encodingConfig, utils.MustCreateGrpcConnection(cfg), forboleXDB),
		staking.NewModule(encodingConfig, utils.MustCreateGrpcConnection(cfg), forboleXDB),
		distribution.NewModule(utils.MustCreateGrpcConnection(cfg), forboleXDB),
		pricefeed.NewModule(encodingConfig, forboleXDB),
	}
}
