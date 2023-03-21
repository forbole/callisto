package resource

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v4/database"

	"github.com/forbole/juno/v4/modules"
	junomessages "github.com/forbole/juno/v4/modules/messages"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/resource module
type Module struct {
	cdc           codec.Codec
	db            *database.Db
	messageParser junomessages.MessageAddressesParser
}

// NewModule returns a new Module instance
func NewModule(
	messageParser junomessages.MessageAddressesParser,
	cdc codec.Codec,
	db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		messageParser: messageParser,
		db:            db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "resource"
}
