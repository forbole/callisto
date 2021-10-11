package slashing

import (
	"github.com/desmos-labs/juno/v2/modules"

	"github.com/forbole/bdjuno/v2/database"
	slashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source"
)

var (
	_ modules.Module      = &Module{}
	_ modules.BlockModule = &Module{}
)

// Module represent x/slashing module
type Module struct {
	db     *database.Db
	source slashingsource.Source

	stakingModule StakingModule
}

// NewModule returns a new Module instance
func NewModule(source slashingsource.Source, stakingModule StakingModule, db *database.Db) *Module {
	return &Module{
		db:            db,
		source:        source,
		stakingModule: stakingModule,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "slashing"
}
