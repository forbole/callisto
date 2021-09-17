package issuer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	issuertypes "github.com/e-money/em-ledger/x/issuer/types"
	"github.com/go-co-op/gocron"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent database/iscn module
type Module struct {
	inflationClient issuertypes.QueryClient
	mintClient      minttypes.QueryClient
	db              *database.Db
}

// NewModule returns a new Module instance
func NewModule(issuerClient issuertypes.QueryClient, mintClient minttypes.QueryClient, db *database.Db) *Module {
	return &Module{
		inflationClient: issuerClient,
		mintClient:      mintClient,
		db:              db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "issuer"
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.mintClient, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, index, msg, m.inflationClient, m.db)
}
