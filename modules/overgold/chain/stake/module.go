package stake

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/database/overgold/chain/stake"
	"github.com/forbole/bdjuno/v4/modules/overgold/chain/stake/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/stake module
type Module struct {
	cdc       codec.Codec
	db        *database.Db
	stakeRepo stake.Repository

	keeper source.Source
}

// NewModule returns a new Module instance
func NewModule(keeper source.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		keeper:    keeper,
		cdc:       cdc,
		db:        db,
		stakeRepo: *stake.NewRepository(db.Sqlx, cdc),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "overgold_stake"
}
