package core

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/core"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/core/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/core module
type Module struct {
	cdc      codec.Codec
	db       *database.Db
	coreRepo core.Repository

	keeper source.Source
}

// NewModule returns a new Module instance
func NewModule(keeper source.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		keeper:   keeper,
		cdc:      cdc,
		db:       db,
		coreRepo: *core.NewRepository(db.Sqlx, cdc),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "overgold_core"
}
