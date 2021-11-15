package distribution

import (
	juno "github.com/forbole/juno/v2/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// Update the validator commissions amount upon reaching interval or if no commission amount is saved in db
	if m.shouldUpdateValidatorsCommissionAmounts(b.Block.Height) {
		go m.updateValidatorsCommissionAmounts(b.Block.Height)
	}

	// Update the delegators commissions amounts upon reaching interval or no rewards saved yet
	if m.shouldUpdateDelegatorRewardsAmounts(b.Block.Height) {
		go m.refreshDelegatorsRewardsAmounts(b.Block.Height)
	}

	return nil
}
