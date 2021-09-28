package authority

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	"github.com/go-co-op/gocron"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represent database/mint module
type Module struct {
	encodingConfig  *params.EncodingConfig
	authorityClient authoritytypes.QueryClient
	db              *database.Db
}

// NewModule returns a new Module instance
func NewModule(
	authorityClient authoritytypes.QueryClient,
	encodingConfig *params.EncodingConfig,
	db *database.Db,
) *Module {
	return &Module{
		encodingConfig:  encodingConfig,
		authorityClient: authorityClient,
		db:              db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "authority"
}

// HandleGenesis implements modules.BlockModule
func (m *Module) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(genesisDoc, appState, m.encodingConfig.Marshaler, m.db)
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.authorityClient, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, msg, m.authorityClient, m.db)
}
