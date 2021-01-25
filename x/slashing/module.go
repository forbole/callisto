package slashing

import (
	"encoding/json"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
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

// Module represent x/slashing module
type Module struct {
	slashingClient slashingtypes.QueryClient
	db             *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(grpcConnection *grpc.ClientConn, db *database.BigDipperDb) *Module {
	return &Module{
		slashingClient: slashingtypes.NewQueryClient(grpcConnection),
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "slashing"
}

// RunAdditionalOperations implements modules.Module
func (m *Module) RunAdditionalOperations() error {
	return nil
}

// RunAsyncOperations implements modules.Module
func (m *Module) RunAsyncOperations() {
}

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(*gocron.Scheduler) error {
	return nil
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(*tmtypes.GenesisDoc, map[string]json.RawMessage) error {
	return nil
}

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(block, m.slashingClient, m.db)
}

// HandleTx implements modules.Module
func (m *Module) HandleTx(*types.Tx) error {
	return nil
}

// HandleMsg implements modules.Module
func (m *Module) HandleMsg(int, sdk.Msg, *types.Tx) error {
	return nil
}
