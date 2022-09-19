package margin

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v3/modules"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/margin/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/margin module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source source.Source
}

// NewModule returns a new Module instance
func NewModule(cdc codec.Codec, source source.Source, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		source: source,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "margin"
}
