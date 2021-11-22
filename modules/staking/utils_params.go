package staking

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"
)

// UpdateParams gets the updated params and stores them inside the database
func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module", "staking").Int64("height", height).
		Msg("updating params")

	params, err := m.source.GetParams(height)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	return m.db.SaveStakingParams(types.NewStakingParams(params, height))
}

// SaveGenesisParams saves the staking parameters at genesis
func (m *Module) SaveGenesisParams(params stakingtypes.Params, height int64) error {
	return m.db.SaveStakingParams(types.NewStakingParams(
		params, height,
	))
}
