package assets

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/database/vipcoin/chain/assets"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/assets module
type Module struct {
	assetRepo assets.Repository
	cdc       codec.Marshaler
	db        *database.Db
	keeper    source.Source
}

// NewModule returns a new Module instance
func NewModule(
	keeper source.Source, cdc codec.Marshaler, db *database.Db,
) *Module {
	return &Module{
		assetRepo: *assets.NewRepository(db.Sqlx, cdc),
		cdc:       cdc,
		db:        db,
		keeper:    keeper,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "overgold_assets"
}
