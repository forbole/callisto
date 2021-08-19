package pricefeed

import (
	"github.com/forbole/bdjuno/types/config"

	"github.com/forbole/bdjuno/database"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/modules"
	"github.com/go-co-op/gocron"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
	_ modules.PeriodicOperationsModule   = &Module{}
)

// Module represents the module that allows to get the token prices
type Module struct {
	cfg            *config.Config
	encodingConfig *params.EncodingConfig
	db             *database.Db
}

// NewModule returns a new Module instance
func NewModule(cfg *config.Config, encodingConfig *params.EncodingConfig, db *database.Db) *Module {
	return &Module{
		cfg:            cfg,
		encodingConfig: encodingConfig,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "pricefeed"
}

// RunAdditionalOperations implements modules.AdditionalOperationsModule
func (m *Module) RunAdditionalOperations() error {
	return StoreTokens(m.cfg, m.db)
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.db)
}
