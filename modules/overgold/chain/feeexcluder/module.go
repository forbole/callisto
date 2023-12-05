package feeexcluder

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/feeexcluder"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/feeexcluder/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/feeexcluder module
type Module struct {
	cdc             codec.Codec
	db              *database.Db
	feeexcluderRepo feeexcluder.Repository

	keeper source.Source
}

// NewModule returns a new Module instance
func NewModule(keeper source.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		keeper:          keeper,
		cdc:             cdc,
		db:              db,
		feeexcluderRepo: *feeexcluder.NewRepository(db.Sqlx, cdc),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "overgold_feeexcluder"
}
