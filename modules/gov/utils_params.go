package gov

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/forbole/callisto/v4/types"
)

// UpdateParams updates the governance parameters for the given height
func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module", "gov").Int64("height", height).
		Msg("updating params")

	params, err := m.source.Params(height)
	if err != nil {
		return fmt.Errorf("error while getting gov params: %s", err)
	}

	return m.db.SaveGovParams(types.NewGovParams(params, height))
}
