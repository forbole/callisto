package mint

import (
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	minttypes "github.com/osmosis-labs/osmosis/x/mint/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module      = &Module{}
	_ modules.BlockModule = &Module{}
)

// Module represent database/mint module
type Module struct {
	mintClient minttypes.QueryClient
	db         *database.Db
}

// NewModule returns a new Module instance
func NewModule(mintClient minttypes.QueryClient, db *database.Db) *Module {
	return &Module{
		mintClient: mintClient,
		db:         db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "mint"
}

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*juno.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(block, m.mintClient, m.db)
}
