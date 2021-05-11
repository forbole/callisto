package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"google.golang.org/grpc"

	"github.com/forbole/bdjuno/modules/common/distribution"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	"github.com/go-co-op/gocron"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var _ modules.Module = &Module{}

// Module represents the x/distr module
type Module struct {
	db          *bigdipperdb.Db
	distrClient distrtypes.QueryClient
}

// NewModule returns a new Module instance
func NewModule(grpConnection *grpc.ClientConn, db *bigdipperdb.Db) *Module {
	return &Module{
		distrClient: distrtypes.NewQueryClient(grpConnection),
		db:          db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "distribution"
}

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	return RegisterPeriodicOps(scheduler, m.distrClient, m.db)
}

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(b *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return distribution.HandleBlock(b, m.distrClient, m.db)
}
