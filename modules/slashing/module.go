package slashing

import (
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/forbole/bdjuno/database"

	"github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	_ modules.Module      = &Module{}
	_ modules.BlockModule = &Module{}
)

// Module represent x/slashing module
type Module struct {
	slashingClient slashingtypes.QueryClient
	db             *database.Db
}

// NewModule returns a new Module instance
func NewModule(slashingClient slashingtypes.QueryClient, db *database.Db) *Module {
	return &Module{
		slashingClient: slashingClient,
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
