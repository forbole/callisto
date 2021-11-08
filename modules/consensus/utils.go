package consensus

import (
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateBlock updates blocks in database
func (m *Module) UpdateBlock(block *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults) error {
	return m.db.UpdateBlockInDatabase(block, blockResults)
}

// UpdateTxs updates txs in database
func (m *Module) UpdateTxs(txHash string, height int64, success bool, messages []string, memo string, signatures []string, signersInfo []byte, fee string, gasWanted int64, gasUsed int64, rawLog string, logs []byte) error {
	return m.db.UpdateTxInDatabase(txHash, height, success, messages, memo, signatures, signersInfo, fee, gasWanted, gasUsed, rawLog, logs)
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
