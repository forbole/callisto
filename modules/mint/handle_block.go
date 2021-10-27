package mint

import (
	juno "github.com/forbole/juno/v2/types"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	return nil
}
