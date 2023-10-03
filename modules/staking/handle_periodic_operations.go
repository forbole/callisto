package staking

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "staking").Msg("setting up periodic tasks")

	// Update the staking pool every 5 mins
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.UpdateStakingPool)
	}); err != nil {
		return fmt.Errorf("error while scheduling staking pool periodic operation: %s", err)
	}

	// refresh proposal validators status snapshots every 5 mins
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.UpdateValidatorStatuses)
	}); err != nil {
		return fmt.Errorf("error while setting up gov period operations: %s", err)
	}

	return nil
}

// UpdateStakingPool reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) UpdateStakingPool() error {
	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}
	log.Debug().Str("module", "staking").Int64("height", block.Height).
		Msg("updating staking pool")

	pool, err := m.GetStakingPool(block.Height)
	if err != nil {
		return fmt.Errorf("error while getting staking pool: %s", err)

	}

	err = m.db.SaveStakingPool(pool)
	if err != nil {
		return fmt.Errorf("error while saving staking pool: %s", err)

	}

	return nil
}
