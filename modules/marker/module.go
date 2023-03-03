package marker

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"github.com/forbole/bdjuno/v3/database"
	markersource "github.com/forbole/bdjuno/v3/modules/marker/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
)

// Module represent database/marker module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source markersource.Source
}

// NewModule returns a new Module instance
func NewModule(source markersource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "marker"
}
