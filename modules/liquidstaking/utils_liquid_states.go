package liquidstaking

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
	"github.com/rs/zerolog/log"
)

// updateLiquidStakingState fetch the liquid staking state and store it into database
func (m *Module) updateLiquidStakingState(height int64) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("getting community pool")

	state, err := m.source.GetLiquidStakingStates(height)
	if err != nil {
		return fmt.Errorf("error while getting liquid staking state: %s", err)
	}

	// Store the pool into the database
	return m.db.SaveLiquidStakingState(types.NewLiquidStakingState(state, height))
}
