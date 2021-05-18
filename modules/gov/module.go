package gov

import (
	"encoding/json"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"google.golang.org/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	encodingConfig *params.EncodingConfig
	govClient      govtypes.QueryClient
	bankClient     banktypes.QueryClient
	db             *database.Db
}

// NewModule returns a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *database.Db) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		govClient:      govtypes.NewQueryClient(grpcConnection),
		bankClient:     banktypes.NewQueryClient(grpcConnection),
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "gov"
}

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(_ *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(appState, m.encodingConfig.Marshaler, m.db)
}

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(_ *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(m.govClient, m.bankClient, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, index, msg, m.govClient, m.encodingConfig.Marshaler, m.db)
}
