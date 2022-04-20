package distribution

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v3/types"
)

// UpdateParams gets the updated params and stores them inside the database
func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module", "distribution").Int64("height", height).
		Msg("updating params")

	params, err := m.source.Params(height)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	return m.db.SaveDistributionParams(types.NewDistributionParams(params, height))

}
