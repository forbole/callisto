package mint

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represent x/mint module
type Module struct {
	mintClient minttypes.QueryClient
	db         *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(grpcConnection *grpc.ClientConn, db *database.BigDipperDb) *Module {
	return &Module{
		mintClient: minttypes.NewQueryClient(grpcConnection),
		db:         db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "mint"
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
	return RegisterPeriodicOps(scheduler, m.mintClient, m.db)
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
