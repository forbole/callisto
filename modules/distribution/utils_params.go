package distribution

import (
	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) UpdateParams(height int64) {
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
