package mint

import (
	"fmt"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"
)

// UpdateParams gets the updated params and stores them inside the database
func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module", "mint").Int64("height", height).
		Msg("updating params")

	params, err := m.source.Params(height)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	return m.db.SaveMintParams(types.NewMintParams(params, height))

}

// SaveGenesisParams saves the governance parameters at genesis
func (m *Module) SaveGenesisParams(params minttypes.Params, height int64) error {
	return m.db.SaveMintParams(types.NewMintParams(
		params, height,
	))
}
