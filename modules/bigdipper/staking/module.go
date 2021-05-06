package staking

import (
	"encoding/json"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/staking module
type Module struct {
	encodingConfig *params.EncodingConfig
	stakingClient  stakingtypes.QueryClient
	bankClient     banktypes.QueryClient
	db             *bigdipperdb.Db
}

// NewModule returns a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *bigdipperdb.Db) *Module {
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
