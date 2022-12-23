package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v3/database"

	"github.com/forbole/bdjuno/v3/modules/auth/source"
	"github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/modules/messages"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/auth module
type Module struct {
	cdc            codec.Codec
	db             *database.Db
	messagesParser messages.MessageAddressesParser
	source         source.Source
}

// NewModule builds a new Module instance
func NewModule(source source.Source, messagesParser messages.MessageAddressesParser, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		messagesParser: messagesParser,
		source:         source,
		cdc:            cdc,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "auth"
}
