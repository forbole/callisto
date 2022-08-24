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

	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateProviders)
	}); err != nil {
		return fmt.Errorf("error while setting up provider periodic operation: %s", err)
	}

	// Setup a cron job to run every 30 minutes
	if _, err := scheduler.Every(30).Minutes().Do(func() {
		utils.WatchMethod(m.updateProviderInventory)
	}); err != nil {
		return err
	}

	return nil
}

// updateProviders fetches from the REST APIs the latest value for the provider list saves it inside the database.
func (m *Module) updateProviders() error {
	log.Debug().
		Str("module", "provider").
		Msg("getting provider list")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get provider addresses
	providers, err := m.getProviderList(height)
	if err != nil {
		return err
	}

	return m.db.SaveProviders(providers, height)
}

// updateProviderInventory fetches from the REST APIs the latest value for the provider resources
// including vCPU, memory and ephemeral storage, and saves them inside the database.
func (m *Module) updateProviderInventory() error {
	log.Debug().
		Str("module", "provider").
		Str("operation", "status").
		Msg("getting provider inventory status")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get provider addresses from the database
	addresses, err := m.db.GetAkashProviders()
	if err != nil {
		return err
	}

	for _, address := range addresses {
		// Get providers inventory status concurrently
		go m.updateProviderInventoryStatus(address, height)
	}

	return nil
}
