package liquidstaking

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
	"github.com/rs/zerolog/log"
)

// updateLiquidStakingState fetch the liquid staking state and store it into database
func (m *Module) updateLiquidStakingState(height int64) error {
	log.Debug().Str("module", "liquidstaking").Int64("height", height).Msg("getting liquid staking state")

	state, err := m.source.GetLiquidStakingStates(height)
	if err != nil {
		return fmt.Errorf("error while getting liquid staking state: %s", err)
	}

	// Store the liquid staking state in the database
	return m.db.SaveLiquidStakingState(types.NewLiquidStakingState(state, height))
}