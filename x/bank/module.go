package bank

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/bank module
type Module struct {
	encodingConfig *params.EncodingConfig
	authClient     authtypes.QueryClient
	bankClient     banktypes.QueryClient
	db             *database.BigDipperDb
}

// NewModule returns a new Module instance
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
	return "bank"
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(appState, m.encodingConfig.Marshaler, m.db)
}

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(block, m.bankClient, m.db)
}

// HandleMsg implements modules.Module
func (m *Module) HandleMsg(_ int, sdkMsg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, sdkMsg, m.authClient, m.bankClient, m.encodingConfig.Marshaler, m.db)
}
