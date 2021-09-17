package staking

import (
	"encoding/json"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/forbole/bdjuno/database"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
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
	cfg            juno.Config
	encodingConfig *params.EncodingConfig
	stakingClient  stakingtypes.QueryClient
	bankClient     banktypes.QueryClient
	distrClient    distrtypes.QueryClient
	db             *database.Db
}

// NewModule returns a new Module instance
func NewModule(
	cfg juno.Config,
	bankClient banktypes.QueryClient, stakingClient stakingtypes.QueryClient, distrClient distrtypes.QueryClient,
	encodingConfig *params.EncodingConfig, db *database.Db,
) *Module {
	return &Module{
		cfg:            cfg,
		encodingConfig: encodingConfig,
		stakingClient:  stakingClient,
		bankClient:     bankClient,
		distrClient:    distrClient,
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
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*juno.Tx, vals *tmctypes.ResultValidators) error {
	return HandleBlock(m.cfg, block, vals, m.stakingClient, m.bankClient, m.distrClient, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, index, msg, m.stakingClient, m.distrClient, m.encodingConfig.Marshaler, m.db)
}
