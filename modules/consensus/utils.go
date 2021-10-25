package consensus

import (
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateBlock updates blocks in database
func (m *Module) UpdateBlock(block *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults) error {
	return m.db.UpdateBlockInDatabase(block, blockResults)
}

// IsBlockMissing checks if block is missing in database
func (m *Module) IsBlockMissing(height int64) bool {
	blockIsMissing := m.db.CheckIfBlockIsMissing(height)
	return blockIsMissing
}

// GetStartingHeight takes starting height value from the config file
func (m *Module) GetStartingHeight() int64 {
	startHeight := m.cfg.StartHeight
	return startHeight
}
