package liquidstaking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	liquidstaking "github.com/forbole/bdjuno/v4/modules/liquidstaking/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/liquidstaking module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source liquidstaking.Source
}

// NewModule returns a new Module instance
func NewModule(source liquidstaking.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "liquidstaking"
}
