package distribution

import (
	"fmt"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"
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

// SaveGenesisParams saves the governance parameters at genesis
func (m *Module) SaveGenesisParams(params distrtypes.Params, height int64) error {
	return m.db.SaveDistributionParams(types.NewDistributionParams(
		params, height,
	))
}
