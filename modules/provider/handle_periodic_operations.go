package provider

import (
	"github.com/forbole/bdjuno/v2/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "mint").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateProviders)
	}); err != nil {
		return err
	}

	return nil
}

// updateProviders fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) updateProviders() error {
	log.Debug().
		Str("module", "provider").
		Str("operation", "update provider").
		Msg("getting provider statuses")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the providers
	providers, err := m.source.GetProviders(height)
	if err != nil {
		return err
	}

	return m.db.SaveProviders(providers, height)
}
