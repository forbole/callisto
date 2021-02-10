package modules

import (
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/modules"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

type Module struct {
	cfg *config.Config
	db  *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(cfg *config.Config, db *database.BigDipperDb) *Module {
	return &Module{
		cfg: cfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "modules"
}

// RunAdditionalOperations implements AdditionalOperationsModule
func (m *Module) RunAdditionalOperations() error {
	return m.db.InsertEnableModules(m.cfg.Cosmos.Modules)
}
