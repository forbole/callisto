package consensus

import (
	"github.com/forbole/bdjuno/v2/database"

	"github.com/forbole/juno/v2/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.BlockModule              = &Module{}
)

// Module implements the consensus utils
type Module struct {
	db            *database.Db
	bankModule    BankModule
	distrModule   DistrModule
	govModule     GovModule
	stakingModule StakingModule
}

// NewModule builds a new Module instance
func NewModule(bankModule BankModule,
	distrModule DistrModule,
	govModule GovModule,
	stakingModule StakingModule,
	db *database.Db) *Module {
	return &Module{
		bankModule:    bankModule,
		distrModule:   distrModule,
		govModule:     govModule,
		stakingModule: stakingModule,
		db:            db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "consensus"
}
