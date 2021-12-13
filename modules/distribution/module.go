package distribution

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v2/types/config"

	distrsource "github.com/forbole/bdjuno/v2/modules/distribution/source"

	"github.com/forbole/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.GenesisModule              = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
	_ modules.PeriodicOperationsModule   = &Module{}
	_ modules.BlockModule                = &Module{}
	_ modules.MessageModule              = &Module{}
)

// Module represents the x/distr module
type Module struct {
	cdc        codec.Codec
	cfg        *Config
	db         *database.Db
	source     distrsource.Source
	bankModule BankModule
}

// NewModule returns a new Module instance
func NewModule(cfg config.Config, source distrsource.Source, bankModule BankModule, cdc codec.Codec, db *database.Db) *Module {
	distrCfg, err := ParseConfig(cfg.GetBytes())
	if err != nil {
		panic(err)
	}

	return &Module{
		cdc:        cdc,
		cfg:        distrCfg,
		db:         db,
		source:     source,
		bankModule: bankModule,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "distribution"
}
