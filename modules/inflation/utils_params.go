package inflation

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/types"
)

// UpdateParams gets the updated params and stores them inside the database
func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module", "inflation").Int64("height", height).
		Msg("updating params")

	// Get inflation params and store in the database
	params, err := m.source.Params(height)
	if err != nil {
		return fmt.Errorf("error while getting inflation params: %s", err)
	}

	return m.db.SaveEvmosInflationParams(types.NewEvmosInflationParams(params, height))
}
