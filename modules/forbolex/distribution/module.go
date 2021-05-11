package distribution

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"google.golang.org/grpc"

	forbolexdb "github.com/forbole/bdjuno/database/forbolex"
	"github.com/forbole/bdjuno/modules/common/distribution"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var _ modules.Module = &Module{}

// Module represents the x/distr module
type Module struct {
	db          *forbolexdb.Db
	distrClient distrtypes.QueryClient
}

// NewModule returns a new Module instance
func NewModule(grpConnection *grpc.ClientConn, db *forbolexdb.Db) *Module {
	return &Module{
		distrClient: distrtypes.NewQueryClient(grpConnection),
		db:          db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "distribution"
}

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(b *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return distribution.HandleBlock(b, m.distrClient, m.db)
}
