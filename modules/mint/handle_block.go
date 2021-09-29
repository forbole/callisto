package mint

import (
	juno "github.com/desmos-labs/juno/v2/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/bdjuno/v2/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*juno.Tx, _ *tmctypes.ResultValidators) error {
	// Update the params
	go m.updateParams(block.Block.Height)

	return nil
}

// updateParams gets the updated params and stores them inside the database
func (m *Module) updateParams(height int64) {
	log.Debug().Str("module", "mint").Int64("height", height).
		Msg("updating params")

	params, err := m.source.Params(height)
	if err != nil {
		log.Error().Str("module", "mint").Err(err).
			Int64("height", height).
			Msg("error while getting params")
		return
	}

	err = m.db.SaveMintParams(types.NewMintParams(params, height))
	if err != nil {
		log.Error().Str("module", "mint").Err(err).
			Int64("height", height).
			Msg("error while saving params")
		return
	}
}
