package auth

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represents the x/auth module
type Module struct {
	encodingConfig *params.EncodingConfig
	authClient     authtypes.QueryClient
	bankClient     banktypes.QueryClient
	db             *database.BigDipperDb
}

// NewModule builds a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *database.BigDipperDb) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		authClient:     authtypes.NewQueryClient(grpcConnection),
		bankClient:     banktypes.NewQueryClient(grpcConnection),
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "auth"
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
	return RegisterOps(scheduler, m.authClient, m.bankClient, m.encodingConfig.Marshaler, m.db)
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return Handler(appState, m.encodingConfig.Marshaler, m.db)
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
