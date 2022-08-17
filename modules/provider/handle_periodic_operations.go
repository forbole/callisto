package provider

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "provider").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(5).Second().Do(func() {
		utils.WatchMethod(m.UpdateInflation)
	}); err != nil {
		return err
	}

	return nil
}

// updateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func (m *Module) UpdateInflation() error {
	log.Debug().
		Str("module", "provider").
		Str("operation", "status").
		Msg("getting provider status")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the inflation
	status, err := m.source.ProviderStatus("akash10fl5f6ukr8kc03mtmf8vckm6kqqwqpc04eruqa", height)
	if err != nil {
		return err
	}

	fmt.Println("status: ", status)

	return nil
}
