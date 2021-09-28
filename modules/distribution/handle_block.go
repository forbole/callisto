package distribution

import (
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(b *tmctypes.ResultBlock, _ []*juno.Tx, _ *tmctypes.ResultValidators) error {
	// Update the params
	go m.updateParams(b.Block.Height)

	// Update the validator commissions
	go m.updateValidatorsCommissionAmounts(b.Block.Height)

	// Update the delegators commissions amounts
	go m.refreshDelegatorsRewardsAmounts(b.Block.Height)

	return nil
}

func (m *Module) updateParams(height int64) {
	log.Debug().Str("module", "distribution").Int64("height", height).
		Msg("updating params")

	params, err := m.source.Params(height)
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Msg("error while getting params")
		return
	}

	err = m.db.SaveDistributionParams(types.NewDistributionParams(params, height))
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).
			Int64("height", height).
			Msg("error while saving params")
		return
	}
}
