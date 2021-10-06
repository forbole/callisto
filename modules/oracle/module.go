package oracle

import (
	"github.com/desmos-labs/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	oraclesource "github.com/forbole/bdjuno/v2/modules/oracle/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/oracle module
type Module struct {
	db     *database.Db
	source oraclesource.Source
}

// NewModule returns a new Module instance
func NewModule(source oraclesource.Source, db *database.Db) *Module {
	return &Module{
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "oracle"
}
