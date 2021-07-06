package mint

import (
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/mint module
type Module struct {
	mintClient minttypes.QueryClient
	db         *database.Db
}

// NewModule returns a new Module instance
func NewModule(mintClient minttypes.QueryClient, db *database.Db) *Module {
	return &Module{
		mintClient: mintClient,
		db:         db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "mint"
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.mintClient, m.db)
}

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*juno.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(block, m.mintClient, m.db)
}
