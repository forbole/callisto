package lpfarm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	lpfarmsource "github.com/forbole/bdjuno/v4/modules/lpfarm/source"
)

var (
	_ modules.Module = &Module{}
)

// Module represent database/lpfarm module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source lpfarmsource.Source
}

// NewModule returns a new Module instance
func NewModule(source lpfarmsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "lpfarm"
}
