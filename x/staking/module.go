package staking

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// Module represents the x/staking module
type Module struct {
	encodingConfig *params.EncodingConfig
	stakingClient  stakingtypes.QueryClient
	db             *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *database.BigDipperDb) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		stakingClient:  stakingtypes.NewQueryClient(grpcConnection),
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "staking"
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
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(appState, m.encodingConfig.Marshaler, m.db)
}

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*types.Tx, vals *tmctypes.ResultValidators) error {
	return HandleBlock(block, vals, m.stakingClient, m.db)
}

// HandleTx implements modules.Module
func (m *Module) HandleTx(*types.Tx) error {
	return nil
}

// HandleMsg implements modules.Module
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, index, msg, m.stakingClient, m.encodingConfig.Marshaler, m.db)
}
