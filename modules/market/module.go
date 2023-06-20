package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/juno/v4/modules"

	marketsource "github.com/forbole/bdjuno/v4/modules/market/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
)

// Module represent x/market module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source marketsource.Source
}

// NewModule returns a new Module instance
func NewModule(source marketsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "market"
}
