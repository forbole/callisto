package distribution

import (
	juno "github.com/desmos-labs/juno/v2/types"

	"github.com/forbole/bdjuno/v2/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(b *tmctypes.ResultBlock, _ []*juno.Tx, _ *tmctypes.ResultValidators) error {
	// Update the params
	go m.updateParams(b.Block.Height)

	// Update the validator commissions
	go m.updateValidatorsCommissionAmounts(b.Block.Height)

	//get interval cfg and verify if func refreshDelegatorsRewardsAmounts() should be executed
	interval := m.cfg.RewardsFrequency
	if interval == 0 {
		log.Debug().Str("module", "distribution").Msg("delegator rewards refresh interval set to 0. Skipping refresh")
		return nil
	}

	hasRewards, err := m.db.HasDelegatorRewards()
	if err != nil {
		log.Error().Str("module", "distribution").Err(err).Int64("height", b.Block.Height).
			Msg("error while checking delegators reward")
	}

	if hasRewards && b.Block.Height%interval != 0 {
		return nil
	}

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
