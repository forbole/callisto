package profiles

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/forbole/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	profilessource "github.com/forbole/bdjuno/v2/modules/profiles/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source profilessource.Source
}

// NewModule returns a new Module instance
func NewModule(cdc codec.Codec, source profilessource.Source, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles"
}
