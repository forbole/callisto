package consensus

import (
	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/juno/v2/types/config"

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
	cfg           *Config
	db            *database.Db
	authModule    AuthModule
	bankModule    BankModule
	distrModule   DistrModule
	govModule     GovModule
	stakingModule StakingModule
}

// NewModule builds a new Module instance
func NewModule(cfg config.Config,
	authModule AuthModule,
	bankModule BankModule,
	distrModule DistrModule,
	govModule GovModule,
	stakingModule StakingModule,
	db *database.Db,
) *Module {
	consCfg, err := ParseConfig(cfg.GetBytes())
	if err != nil {
		panic(err)
	}
	return &Module{
		cfg: consCfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "consensus"
}
