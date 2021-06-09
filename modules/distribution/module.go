package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"google.golang.org/grpc"

	"github.com/forbole/bdjuno/database"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represents the x/distr module
type Module struct {
	db          *database.Db
	distrClient distrtypes.QueryClient
}

// NewModule returns a new Module instance
func NewModule(grpConnection *grpc.ClientConn, db *database.Db) *Module {
	return &Module{
		distrClient: distrtypes.NewQueryClient(grpConnection),
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
	return HandleBlock(b, m.distrClient, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *types.Tx) error {
	return HandleMsg(tx, msg, m.distrClient, m.db)
}
