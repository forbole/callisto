package consensus

import (
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/juno/v2/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateBlock updates blocks in database
func (m *Module) UpdateBlock(block *tmctypes.ResultBlock) error {
	return m.db.UpdateBlockInDatabase(block)
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

// HandleMessages updates messages in database
func (m *Module) HandleMessages(txDetails *types.Tx) error {

	// Handle messages
	for index, msg := range txDetails.GetMsgs() {
		switch msg.Route() {
		case banktypes.ModuleName:
			messageErr := m.bankModule.HandleMsg(index, msg, txDetails)
			if messageErr != nil {
				return fmt.Errorf("error when updatig bank module Handle Message: %s", messageErr)
			}
		case govtypes.ModuleName:
			messageErr := m.govModule.HandleMsg(index, msg, txDetails)
			if messageErr != nil {
				return fmt.Errorf("error when updatig gov module Handle Message: %s", messageErr)
			}
		case stakingtypes.ModuleName:
			messageErr := m.stakingModule.HandleMsg(index, msg, txDetails)
			if messageErr != nil {
				return fmt.Errorf("error when updatig staking module Handle Message: %s", messageErr)
			}
		case distrtypes.ModuleName:
			messageErr := m.distrModule.HandleMsg(index, msg, txDetails)
			if messageErr != nil {
				return fmt.Errorf("error when updatig distr module Handle Message: %s", messageErr)
			}
		}

		err := m.UpdateTxs(index, txDetails)
		if err != nil {
			return fmt.Errorf("error when updatig transactions tx %s", err)
		}
	}
	return nil

}
