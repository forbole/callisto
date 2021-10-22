package consensus

import (
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateBlock updated blocks in database
func (m *Module) UpdateBlock(block *tmctypes.ResultBlock) error {
	return m.db.UpdateBlockInDatabase(block)
}
