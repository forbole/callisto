package modules

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

type Module struct {
	cfg *config.Config
	db  *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(cfg *config.Config, db *database.BigDipperDb) *Module {
	return &Module{
		cfg: cfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "modules"
}

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(*gocron.Scheduler) error {
	return nil
}

// RunAsyncOperations implements modules.Module
func (m *Module) RunAsyncOperations() {
}

// RunAdditionalOperations implements modules.Module
func (m *Module) RunAdditionalOperations() error {
	return m.db.InsertEnableModules(m.cfg.CosmosConfig.Modules)
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(*tmtypes.GenesisDoc, map[string]json.RawMessage) error {
	return nil
}

// HandleBlock implements modules.Module
func (m Module) HandleBlock(*tmctypes.ResultBlock, []*types.Tx, *tmctypes.ResultValidators) error {
	return nil
}

// HandleTx implements modules.Module
func (m *Module) HandleTx(*types.Tx) error {
	return nil
}

// HandleMsg implements modules.Module
func (m *Module) HandleMsg(int, sdk.Msg, *types.Tx) error {
	return nil
}
