package message_type

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v6/modules"
	"github.com/forbole/juno/v6/modules/messages"

	"github.com/forbole/callisto/v4/database"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

type Module struct {
	cdc    codec.Codec
	db     *database.Db
	parser messages.MessageAddressesParser
}

func NewModule(parser messages.MessageAddressesParser, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		parser: parser,
		cdc:    cdc,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "message_type"
}
