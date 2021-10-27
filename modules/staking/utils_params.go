package staking

import (
	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"
)

// UpdateParams gets the updated params and stores them inside the database
func (m *Module) UpdateParams(height int64) {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating params")

	params, err := m.source.GetParams(height)
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while getting params")
		return
	}

	err = m.db.SaveStakingParams(types.NewStakingParams(params, height))
	if err != nil {
		log.Error().Str("module", "staking").Err(err).
			Int64("height", height).
			Msg("error while saving params")
		return
	}
}
