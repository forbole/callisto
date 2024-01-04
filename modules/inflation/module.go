package inflation

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/juno/v4/modules"

	"github.com/forbole/bdjuno/v4/database"
	inflationsource "github.com/forbole/bdjuno/v4/modules/inflation/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/mint module
type Module struct {
	cdc    codec.Marshaler
	db     *database.Db
	source inflationsource.Source
}

// NewModule returns a new Module instance
func NewModule(cdc codec.Marshaler, source inflationsource.Source, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "inflation"
}
