package distribution

import (
	"github.com/cosmos/cosmos-sdk/codec"

	distrsource "github.com/forbole/callisto/v4/modules/distribution/source"

	"github.com/forbole/juno/v5/modules"

	"github.com/forbole/callisto/v4/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represents the x/distr module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source distrsource.Source
}

// NewModule returns a new Module instance
func NewModule(source distrsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "distribution"
}
