/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package wallets

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v2/modules"
	junomessages "github.com/forbole/juno/v2/modules/messages"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/database/vipcoin/chain/wallets"
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represents the x/wallets module
type Module struct {
	cdc           codec.Marshaler
	db            *database.Db
	walletsRepo   wallets.Repository
	messageParser junomessages.MessageAddressesParser
	keeper        source.Source
}

// NewModule returns a new Module instance
func NewModule(
	messageParser junomessages.MessageAddressesParser, keeper source.Source, cdc codec.Marshaler, db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		db:            db,
		walletsRepo:   *wallets.NewRepository(db.Sqlx, cdc),
		messageParser: messageParser,
		keeper:        keeper,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "vipcoin_wallets"
}
