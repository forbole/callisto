package staking

import (
	"encoding/json"

	"github.com/forbole/bdjuno/modules/common/staking"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	forbolexdb "github.com/forbole/bdjuno/database/forbolex"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// Module represents the x/staking module
type Module struct {
	encodingConfig *params.EncodingConfig
	stakingClient  stakingtypes.QueryClient
	db             *forbolexdb.Db
}

// NewModule returns a new Module instance
func NewModule(encodingConfig *params.EncodingConfig, grpcConnection *grpc.ClientConn, db *forbolexdb.Db) *Module {
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

// HandleGenesis implements  modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return staking.HandleGenesis(doc, appState, m.encodingConfig.Marshaler, m.db)
}

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(block.Block, m.stakingClient, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, index, msg, m.stakingClient, m.encodingConfig.Marshaler, m.db)
}
