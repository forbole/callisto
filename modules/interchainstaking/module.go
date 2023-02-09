package interchainstaking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"github.com/forbole/bdjuno/v4/database"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
)

// Module represent database/interchainstaking module
type Module struct {
	cdc codec.Codec
	db  *database.Db
}

// NewModule returns a new Module instance
func NewModule(cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc: cdc,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "interchainstaking"
}
