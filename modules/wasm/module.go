package wasm

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/bdjuno/v3/database"
	wasmsource "github.com/forbole/bdjuno/v3/modules/wasm/source"

	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/feegrant module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source wasmsource.Source
}

// NewModule returns a new Module instance
func NewModule(source wasmsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "wasm"
}
