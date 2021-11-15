package consensus

import (
	"github.com/forbole/juno/v2/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateBlock updates blocks in database
func (m *Module) UpdateBlock(block *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults) error {
	return m.db.UpdateBlockInDatabase(block, blockResults)
}

// UpdateTxs updates txs in database
func (m *Module) UpdateTxs(i int, tx *types.Tx) error {
	return m.db.UpdateTxInDatabase(i, tx)
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
