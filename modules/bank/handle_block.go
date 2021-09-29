package bank

import (
	"github.com/desmos-labs/juno/v2/types"

	"github.com/rs/zerolog/log"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*types.Tx, _ *tmctypes.ResultValidators) error {
	err := m.updateSupply(block.Block.Height)
	if err != nil {
		log.Error().Str("module", "bank").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating supply")
	}

	return nil
}

// updateSupply updates the supply of all the tokens for the given height
func (m *Module) updateSupply(height int64) error {
	log.Debug().Str("module", "bank").Int64("height", height).Msg("updating supply")

	supply, err := m.keeper.GetSupply(height)
	if err != nil {
		return err
	}

	return m.db.SaveSupply(supply, height)
}
