package provider

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/modules/utils"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "provider").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.updateProviders)
	}); err != nil {
		return fmt.Errorf("error while setting up provider periodic operation: %s", err)
	}

	return nil
}

// updateProviders fetches from the REST APIs the latest value for the provider list saves it inside the database.
func (m *Module) updateProviders() error {
	log.Debug().
		Str("module", "provider").
		Msg("updating provider list and inventory status")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get providers info
	providers, err := m.getProviderList(height)
	if err != nil {
		return err
	}

	// Save providers into the database
	err = m.db.SaveProviders(providers, height)
	if err != nil {
		return err
	}

	// Get provider inventory statuses concurrently and save them into the database
	for _, p := range providers {
		go m.updateProviderInventoryStatus(p.OwnerAddress, p.HostURI, height)
	}

	return nil
}
