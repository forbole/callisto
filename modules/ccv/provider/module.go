package provider

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/bdjuno/v4/database"
	providersource "github.com/forbole/bdjuno/v4/modules/ccv/provider/source"
)

var (
	_ modules.Module = &Module{}
)

// Module represent database/ccv/provider module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source providersource.Source
}

// NewModule returns a new Module instance
func NewModule(source providersource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "ccvprovider"
}
