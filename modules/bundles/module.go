package bundles

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"github.com/forbole/bdjuno/v4/database"
	bundlessource "github.com/forbole/bdjuno/v4/modules/bundles/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent database/bundles module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source bundlessource.Source
}

// NewModule returns a new Module instance
func NewModule(source bundlessource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "bundles"
}
