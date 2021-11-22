package history

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/modules/messages"
	"github.com/forbole/juno/v2/types/config"

	"github.com/forbole/bdjuno/v2/database"
)

const (
	moduleName = "history"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represents the module that allows to store historic information
type Module struct {
	cfg config.ChainConfig
	cdc codec.Codec
	db  *database.Db

	getAddresses messages.MessageAddressesParser
}

// NewModule allows to build a new Module instance
func NewModule(cfg config.ChainConfig, messagesParser messages.MessageAddressesParser, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cfg:          cfg,
		cdc:          cdc,
		db:           db,
		getAddresses: messagesParser,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return moduleName
}
