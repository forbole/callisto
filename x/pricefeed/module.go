package pricefeed

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/modules"
	"github.com/go-co-op/gocron"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represent x/pricefeed module
type Module struct {
	encodingConfig *params.EncodingConfig
	db             *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, db *database.BigDipperDb) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "pricefeed"
}

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(doc, appState, m.encodingConfig.Marshaler, m.db)
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.db)
}
