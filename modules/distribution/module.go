package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types/config"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represents the x/distr module
type Module struct {
	cfg         *config.Config
	db          *database.Db
	bankClient  banktypes.QueryClient
	distrClient distrtypes.QueryClient
}

// NewModule returns a new Module instance
func NewModule(cfg *config.Config, bankClient banktypes.QueryClient, distrClient distrtypes.QueryClient, db *database.Db) *Module {
	return &Module{
		cfg:         cfg,
		bankClient:  bankClient,
		distrClient: distrClient,
		db:          db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "distribution"
}

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.distrClient, m.db)
}

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(b *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(m.cfg, b, m.distrClient, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, msg, m.distrClient, m.bankClient, m.db)
}
