package accounts

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/database/overgold/chain/accounts"
	"github.com/forbole/bdjuno/v2/database/overgold/chain/assets"
	"github.com/forbole/bdjuno/v2/database/overgold/chain/wallets"
	"github.com/forbole/bdjuno/v2/modules/overgold/chain/accounts/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/accounts module
type Module struct {
	cdc         codec.Marshaler
	db          *database.Db
	walletsRepo wallets.Repository
	accountRepo accounts.Repository
	assetRepo   assets.Repository

	keeper source.Source
}

// NewModule returns a new Module instance
func NewModule(
	keeper source.Source, cdc codec.Marshaler, db *database.Db,
) *Module {
	return &Module{
		cdc:         cdc,
		db:          db,
		accountRepo: *accounts.NewRepository(db.Sqlx, cdc),
		walletsRepo: *wallets.NewRepository(db.Sqlx, cdc),
		assetRepo:   *assets.NewRepository(db.Sqlx, cdc),
		keeper:      keeper,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "overgold_accounts"
}
