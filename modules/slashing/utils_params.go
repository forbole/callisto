package slashing

import (
	"fmt"

	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// UpdateParams gets the slashing params for the given height, and stores them inside the database
func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module", "slashing").Int64("height", height).Msg("updating params")

	params, err := m.source.GetParams(height)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	return m.db.SaveSlashingParams(types.NewSlashingParams(params, height))
}

// SaveGenesisParams saves the governance parameters at genesis
func (m *Module) SaveGenesisParams(params slashingtypes.Params, height int64) error {
	return m.db.SaveSlashingParams(types.NewSlashingParams(
		params, height,
	))
}
