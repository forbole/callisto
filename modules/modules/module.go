package modules

import (
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

type Module struct {
	cfg juno.Config
	db  *database.Db
}

// NewModule returns a new Module instance
func NewModule(cfg juno.Config, db *database.Db) *Module {
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
	return m.db.InsertEnableModules(m.cfg.GetCosmosConfig().GetModules())
}
