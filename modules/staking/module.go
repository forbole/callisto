package staking

import (
	"encoding/json"

	"github.com/forbole/bdjuno/database"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
	db             *database.Db
}

// NewModule returns a new Module instance
func NewModule(
	bankClient banktypes.QueryClient, stakingClient stakingtypes.QueryClient,
	encodingConfig *params.EncodingConfig, db *database.Db,
) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		stakingClient:  stakingClient,
		bankClient:     bankClient,
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
