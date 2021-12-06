package slashing

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	slashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
)

// Module represent x/slashing module
type Module struct {
	cdc    codec.Marshaler
	db     *database.Db
	source slashingsource.Source

	stakingModule StakingModule
}

// NewModule returns a new Module instance
func NewModule(source slashingsource.Source, stakingModule StakingModule, cdc codec.Marshaler, db *database.Db) *Module {
	return &Module{
		cdc:           cdc,
		db:            db,
		source:        source,
		stakingModule: stakingModule,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "slashing"
}
