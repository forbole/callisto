package consensus

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module implements the consensus operations
type Module struct {
	cp *client.Proxy
	db *database.BigDipperDb
}

// NewModule builds a new Module instance
func NewModule(cp *client.Proxy, db *database.BigDipperDb) *Module {
	return &Module{
		cp: cp,
		db: db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "consensus"
}

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return Register(scheduler, m.db)
}

// RunAdditionalOperations implements modules.Module
func (m *Module) RunAdditionalOperations() error {
	return nil
}

// RunAsyncOperations implements modules.Module
func (m *Module) RunAsyncOperations() {
	go ListenOperation(m.cp, m.db)
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, _ map[string]json.RawMessage) error {
	return HandleGenesis(doc, m.db)
}

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(b *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(b, m.db)
}

// HandleTx implements modules.Module
func (m *Module) HandleTx(*types.Tx) error {
	return nil
}

// HandleMsg implements modules.Module
func (m *Module) HandleMsg(int, sdk.Msg, *types.Tx) error {
	return nil
}
