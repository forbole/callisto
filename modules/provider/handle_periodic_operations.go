package provider

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/ovrclk/akash/provider"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "provider").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.updateProviders)
	}); err != nil {
		return fmt.Errorf("error while setting up consensus periodic operation: %s", err)
	}

	// Setup a cron job to run every 10 minutes
	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.updateProviderProvisionStatus)
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
	providers, err := m.source.GetProviders(height)
	if err != nil {
		return err
	}

	return m.db.SaveProviders(providers, height)
}

// updateProviderProvisionStatus fetches from the REST APIs the latest value for the provider resources
// including vCPU, memory and ephemeral storage, and saves them inside the database.
func (m *Module) updateProviderProvisionStatus() error {
	log.Debug().
		Str("module", "provider").
		Str("operation", "provider status").
		Msg("getting provider provision")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get provider addresses
	addresses, err := m.db.GetAkashProviders()
	if err != nil {
		return err
	}

	// Get provider statuses
	statuses, err := m.getProviderProvisionStatuses(addresses)
	if err != nil {
		return err
	}

	return m.db.SaveProviderProvisionStatuses(statuses, height)
}

func (m *Module) getProviderProvisionStatuses(addresses []string) ([]*provider.Status, error) {

	providerStatuses := make([]*provider.Status, len(addresses))
	for i, address := range addresses {
		// Get the provision status of a provider
		status, err := m.source.GetProviderProvisionStatus(address)
		if err != nil {
			return nil, err
		}
		providerStatuses[i] = status
	}

	return providerStatuses, nil
}
