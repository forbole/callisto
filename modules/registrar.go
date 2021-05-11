package modules

import (
	"fmt"

	"github.com/forbole/bdjuno/types/config"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/registrar"
	juno "github.com/desmos-labs/juno/types"

	bigdippermodules "github.com/forbole/bdjuno/modules/bigdipper"
	forbolexmodules "github.com/forbole/bdjuno/modules/forbolex"
)

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents a modules registrar that decides what modules to register based on the application type
type Registrar struct{}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar() *Registrar {
	return &Registrar{}
}

// BuildModules implements registrar.Registrar
func (r *Registrar) BuildModules(
	junoCfg juno.Config, encodingConfig *params.EncodingConfig, sdkConfig *sdk.Config, db db.Database, cp *client.Proxy,
) modules.Modules {
	cfg, ok := junoCfg.(*config.Config)
	if !ok {
		panic(fmt.Errorf("invalid configuration type: %T", junoCfg))
	}

	var reg registrar.Registrar
	switch cfg.GetDataType() {
	case config.DataTypeUpdated:
		reg = bigdippermodules.NewRegistrar()

	case config.DataTypeHistoric:
		reg = forbolexmodules.NewRegistrar()

	default:
		panic(fmt.Errorf("invalid application type: %s", cfg.GetDataType()))
	}

	return reg.BuildModules(junoCfg, encodingConfig, sdkConfig, db, cp)
}
