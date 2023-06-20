package market

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "market").Msg("setting up periodic tasks")

	// Update the market leases state every 30 minutes
	if _, err := scheduler.Every(30).Minutes().Do(func() {
		utils.WatchMethod(m.getLatestActiveLeases)
	}); err != nil {
		return fmt.Errorf("error while scheduling market periodic operation: %s", err)
	}

	return nil
}

// getLatestActiveLeases gets the latest status of all active leases from the chain and stores inside the database
func (m *Module) getLatestActiveLeases() error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	return m.updateLeases(height)
}
