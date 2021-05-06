package staking

import (
	"encoding/json"

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

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// TODO: Add genesis handling
	return nil
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, index, msg, m.stakingClient, m.db)
}
