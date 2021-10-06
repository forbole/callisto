package slashing

import (
	juno "github.com/desmos-labs/juno/v2/types"

	"github.com/forbole/bdjuno/v2/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(block *tmctypes.ResultBlock, _ []*juno.Tx, _ *tmctypes.ResultValidators) error {
	// Update the signing infos
	err := m.updateSigningInfo(block.Block.Height)
	if err != nil {
		log.Error().Str("module", "slashing").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating signing info")
	}

	err = m.updateSlashingParams(block.Block.Height)
	if err != nil {
		log.Error().Str("module", "slashing").Int64("height", block.Block.Height).
			Err(err).Msg("error while updating params")
	}

	return nil
}

// updateSigningInfo reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateSigningInfo(height int64) error {
	log.Debug().Str("module", "slashing").Int64("height", height).
		Msg("updating signing info")

	signingInfos, err := m.getSigningInfos(height)
	if err != nil {
		return err
	}

	return m.db.SaveValidatorsSigningInfos(signingInfos)
}

// updateSlashingParams gets the slashing params for the given height, and stores them inside the database
func (m *Module) updateSlashingParams(height int64) error {
	log.Debug().Str("module", "slashing").Int64("height", height).
		Msg("updating slashing params")

	params, err := m.source.GetParams(height)
	if err != nil {
		return err
	}

	return m.db.SaveSlashingParams(types.NewSlashingParams(params, height))
}
