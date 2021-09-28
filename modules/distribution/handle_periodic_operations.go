package distribution

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "distribution").Msg("setting up periodic tasks")

	// Update the community pool every 1 hour
	if _, err := scheduler.Every(1).Hour().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return m.getLatestCommunityPool() })
	}); err != nil {
		return fmt.Errorf("error while scheduling distribution peridic operation: %s", err)
	}

	return nil
}

// getLatestCommunityPool gets the latest community pool from the chain and stores inside the database
func (m *Module) getLatestCommunityPool() error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	return m.updateCommunityPool(height)
}
