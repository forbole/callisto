package distribution

import (
	"encoding/json"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"google.golang.org/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represents the x/distr module
type Module struct {
	db          *database.BigDipperDb
	distrClient distrtypes.QueryClient
}

// NewModule returns a new Module instance
func NewModule(grpConnection *grpc.ClientConn, db *database.BigDipperDb) *Module {
	return &Module{
		distrClient: distrtypes.NewQueryClient(grpConnection),
		db:          db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "distribution"
}

// RunAdditionalOperations implements modules.Module
func (m *Module) RunAdditionalOperations() error {
	return nil
}

// RunAsyncOperations implements modules.Module
func (m *Module) RunAsyncOperations() {
}

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.distrClient, m.db)
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(*tmtypes.GenesisDoc, map[string]json.RawMessage) error {
	return nil
}

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(*tmctypes.ResultBlock, []*types.Tx, *tmctypes.ResultValidators) error {
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
