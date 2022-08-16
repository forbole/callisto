package market

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v3/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "market").Msg("setting up periodic tasks")

	// Update the community pool every 30 minutes
	if _, err := scheduler.Every(30).Minutes().Do(func() {
		utils.WatchMethod(m.getLatestLeasesStatus)
	}); err != nil {
		return fmt.Errorf("error while scheduling market periodic operation: %s", err)
	}

	return nil
}

// getLatestLeasesStatus gets the latest status of all leases from the chain and stores inside the database
func (m *Module) getLatestLeasesStatus() error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	return m.updateLeases(height)
}
