package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/bdjuno/v2/database"
	marketsource "github.com/forbole/bdjuno/v2/modules/market/source"
	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represents the x/provider module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source marketsource.Source
}

// NewModule returns a new Module instance
func NewModule(
	source marketsource.Source, cdc codec.Codec, db *database.Db,
) *Module {
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
