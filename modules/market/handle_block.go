package market

import (
	"fmt"

	juno "github.com/forbole/juno/v2/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, results *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// Update lease statuses
	err := m.updateLeases(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating leases: %s", err)
	}

	return nil
}

// updateLeases queries the leases and stores them inside the database
func (m *Module) updateLeases(height int64) error {
	log.Debug().Str("module", "market").
		Int64("height", height).Msg("updating leases")

	leases, err := m.getLeases(height)
	if err != nil {
		return err
	}

	return m.db.SaveLeases(leases, height)
}
