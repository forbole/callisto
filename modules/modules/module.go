package modules

import (
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types/config"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

type Module struct {
	cfg config.ChainConfig
	db  *database.Db
}

// NewModule returns a new Module instance
func NewModule(cfg config.ChainConfig, db *database.Db) *Module {
	return &Module{
		cfg: cfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "modules"
}
