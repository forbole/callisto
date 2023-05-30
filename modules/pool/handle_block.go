package pool

import (
	"fmt"

	juno "github.com/forbole/juno/v4/types"
	"github.com/rs/zerolog/log"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	log.Debug().Str("module", "pool").Int64("height", block.Block.Height).
		Msg("updating pools")

	// Update the protocol validators
	err := m.UpdatePools(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating pools: %s", err)
	}

	return nil
}
