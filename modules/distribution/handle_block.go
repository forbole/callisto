package distribution

import (
	juno "github.com/forbole/juno/v2/types"

	"github.com/forbole/bdjuno/v2/types"

	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// Update the params
	go m.updateParams(b.Block.Height)

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
