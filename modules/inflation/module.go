package inflation

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/modules"
	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
	"github.com/go-co-op/gocron"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/mint module
type Module struct {
	encodingConfig  *params.EncodingConfig
	inflationClient inflationtypes.QueryClient
	db              *database.Db
}

// NewModule returns a new Module instance
func NewModule(
	inflationClient inflationtypes.QueryClient,
	encodingConfig *params.EncodingConfig,
	db *database.Db,
) *Module {
	return &Module{
		encodingConfig:  encodingConfig,
		inflationClient: inflationClient,
		db:              db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "inflation"
}

// HandleBlock implements modules.BlockModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(doc, appState, m.encodingConfig.Marshaler, m.db)
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.inflationClient, m.db)
}
