package stakeibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"github.com/forbole/bdjuno/v4/database"
	stakeibcsource "github.com/forbole/bdjuno/v4/modules/stakeibc/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
)

// Module represent database/stakeibc module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source stakeibcsource.Source
}

// NewModule returns a new Module instance
func NewModule(source stakeibcsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "stakeibc"
}
