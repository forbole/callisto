package pricefeed

import (
	"encoding/json"

	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/modules"
	"github.com/go-co-op/gocron"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.GenesisModule            = &Module{}
)

// Module represents the module that allows to get the token prices
type Module struct {
	cfg            juno.Config
	encodingConfig *params.EncodingConfig
	db             *database.Db
}

// NewModule returns a new Module instance
func NewModule(cfg juno.Config, encodingConfig *params.EncodingConfig, db *database.Db) *Module {
	return &Module{
		cfg:            cfg,
		encodingConfig: encodingConfig,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "pricefeed"
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.db)
}

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(m.cfg, doc, appState, m.encodingConfig.Marshaler, m.db)
}
