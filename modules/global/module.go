package global

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"github.com/forbole/bdjuno/v4/database"
	globalsource "github.com/forbole/bdjuno/v4/modules/global/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent database/global module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source globalsource.Source
}

// NewModule returns a new Module instance
func NewModule(source globalsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "global"
}
