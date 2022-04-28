package banking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/database/vipcoin/chain/accounts"
	"github.com/forbole/bdjuno/v2/database/vipcoin/chain/assets"
	"github.com/forbole/bdjuno/v2/database/vipcoin/chain/banking"
	"github.com/forbole/bdjuno/v2/database/vipcoin/chain/wallets"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/banking module
type Module struct {
	cdc          codec.Marshaler
	db           *database.Db
	bankingRepo  banking.Repository
	walletsRepo  wallets.Repository
	assetRepo    assets.Repository
	accountsRepo accounts.Repository
	keeper       source.Source
	storeKey     sdk.StoreKey
}

// NewModule returns a new Module instance
func NewModule(
	keeper source.Source, storeKey sdk.StoreKey, cdc codec.Marshaler, db *database.Db) *Module {
	return &Module{
		cdc:          cdc,
		db:           db,
		bankingRepo:  *banking.NewRepository(db.Sqlx, cdc),
		walletsRepo:  *wallets.NewRepository(db.Sqlx, cdc),
		assetRepo:    *assets.NewRepository(db.Sqlx, cdc),
		accountsRepo: *accounts.NewRepository(db.Sqlx, cdc),
		keeper:       keeper,
		storeKey:     storeKey,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "vipcoin_banking"
}