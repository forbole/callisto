package slashing

import (
	slashingsource "github.com/forbole/bdjuno/modules/slashing/source"

	"github.com/forbole/bdjuno/database"

	"github.com/desmos-labs/juno/v2/modules"
)

var (
	_ modules.Module      = &Module{}
	_ modules.BlockModule = &Module{}
)

// Module represent x/slashing module
type Module struct {
	db     *database.Db
	source slashingsource.Source
}

// NewModule returns a new Module instance
func NewModule(source slashingsource.Source, db *database.Db) *Module {
	return &Module{
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "slashing"
}
