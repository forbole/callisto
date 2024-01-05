package oracle

import (
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	oraclesource "github.com/forbole/bdjuno/v4/modules/oracle/source"
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
