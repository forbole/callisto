package shield

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v3/database"

	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/bank module
type Module struct {
	cdc codec.Codec
	db  *database.Db

	// source source.Source
}

// NewModule returns a new Module instance
func NewModule(
	// source source.Source,
	cdc codec.Codec, db *database.Db,
) *Module {
	return &Module{
		cdc: cdc,
		db:  db,
		// source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "shield"
}
