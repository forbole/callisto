package wormhole

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	wormholesource "github.com/forbole/bdjuno/v4/modules/wormhole/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represents the x/wormhole module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source wormholesource.Source
}

// NewModule returns a new Module instance
func NewModule(
	source wormholesource.Source, cdc codec.Codec, db *database.Db,
) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "wormhole"
}
