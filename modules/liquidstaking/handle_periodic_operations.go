package liquidstaking

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v3/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "liquidstaking").Msg("setting up periodic tasks")

	// Update the liquid staking state every 1 day
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.getLatestLiquidStakingState)
	}); err != nil {
		return fmt.Errorf("error while scheduling liquidstaking peridic operation: %s", err)
	}

	return nil
}

// getLatestLiquidStakingState gets the latest liquid staking state from the chain and stores inside the database
func (m *Module) getLatestLiquidStakingState() error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	return m.updateLiquidStakingState(height)
}
