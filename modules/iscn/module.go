package iscn

import (
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.BlockModule              = &Module{}
)

// Module represent database/iscn module
type Module struct {
	iscnClient iscntypes.QueryClient
	db         *database.Db
}

// NewModule returns a new Module instance
func NewModule(iscnClient iscntypes.QueryClient, db *database.Db) *Module {
	return &Module{
		iscnClient: iscnClient,
		db:         db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "iscn"
}


// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*juno.Tx, _ *tmctypes.ResultValidators) error {
	return HandleBlock(block, m.iscnClient, m.db)
}
