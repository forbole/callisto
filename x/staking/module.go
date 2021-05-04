package staking

import (
	"encoding/json"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/go-co-op/gocron"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represents the x/staking module
type Module struct {
	encodingConfig *params.EncodingConfig
	stakingClient  stakingtypes.QueryClient
	bankClient     banktypes.QueryClient
	db             *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *database.BigDipperDb) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		stakingClient:  stakingtypes.NewQueryClient(grpcConnection),
		bankClient:     banktypes.NewQueryClient(grpcConnection),
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "staking"
}

// RegisterPeriodicOperations implements PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.stakingClient, m.encodingConfig.Marshaler, m.db)
}

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(doc, appState, m.encodingConfig.Marshaler, m.db)
}

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*types.Tx, vals *tmctypes.ResultValidators) error {
	return HandleBlock(block, vals, m.stakingClient, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, index, msg, m.stakingClient, m.bankClient, m.encodingConfig.Marshaler, m.db)
}
