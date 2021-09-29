package mint

import (
	"github.com/desmos-labs/juno/v2/modules"

	"github.com/forbole/bdjuno/database"
	mintsource "github.com/forbole/bdjuno/modules/mint/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/mint module
type Module struct {
	db     *database.Db
	source mintsource.Source
}

// NewModule returns a new Module instance
func NewModule(source mintsource.Source, db *database.Db) *Module {
	return &Module{
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "mint"
}
