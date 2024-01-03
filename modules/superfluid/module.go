package superfluid

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	superfluidsource "github.com/forbole/bdjuno/v4/modules/superfluid/source"
)

var (
	_ modules.Module = &Module{}
)

// Module represent database/superfluid module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source superfluidsource.Source
}

// NewModule returns a new Module instance
func NewModule(source superfluidsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "superfluid"
}
