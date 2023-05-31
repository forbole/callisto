package pool

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "pool").Msg("setting up periodic tasks")

	// Update the pool every 5 mins
	if _, err := scheduler.Every(5).Minutes().Do(func() {
		utils.WatchMethod(m.UpdatePoolList)
	}); err != nil {
		return fmt.Errorf("error while scheduling pool periodic operation: %s", err)
	}

	return nil
}

// UpdatePoolList queries current pool list and stores it inside the database
func (m *Module) UpdatePoolList() error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	log.Debug().Str("module", "pool").Int64("height", height).
		Msg("updating pool list")

	err = m.UpdatePools(height)
	if err != nil {
		return fmt.Errorf("error while saving pool list: %s", err)

	}

	return nil
}
