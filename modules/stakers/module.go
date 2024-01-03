package stakers

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	stakerssource "github.com/forbole/bdjuno/v4/modules/stakers/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent database/stakers module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source stakerssource.Source
}

// NewModule returns a new Module instance
func NewModule(source stakerssource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "stakers"
}
