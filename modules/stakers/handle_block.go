package stakers

import (
	"fmt"

	juno "github.com/forbole/juno/v5/types"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {

	// Update the protocol validators
	err := m.UpdateProtocolValidatorsInfo(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating protocol validators: %s", err)
	}

	return nil
}
