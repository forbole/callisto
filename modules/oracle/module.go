package oracle

import (
	oracletypes "github.com/bandprotocol/chain/v2/x/oracle/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/go-co-op/gocron"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/oracle module
type Module struct {
	oracleClient oracletypes.QueryClient
	db           *database.Db
}

// NewModule returns a new Module instance
func NewModule(oracleClient oracletypes.QueryClient, db *database.Db) *Module {
	return &Module{
		oracleClient: oracleClient,
		db:           db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "oracle"
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.oracleClient, m.db)
}
