package slashing

import (
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"google.golang.org/grpc"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
)

var _ modules.Module = &Module{}

// Module represent x/slashing module
type Module struct {
	slashingClient slashingtypes.QueryClient
	db             *database.BigDipperDb
}

// NewModule returns a new Module instance
func NewModule(grpcConnection *grpc.ClientConn, db *database.BigDipperDb) *Module {
	return &Module{
		slashingClient: slashingtypes.NewQueryClient(grpcConnection),
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "slashing"
}

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(block, m.slashingClient, m.db)
}
