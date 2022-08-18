package provider

import (
	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/ovrclk/akash/provider"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "provider").Msg("setting up periodic tasks")

	// Setup a cron job to run every 10 minutes
	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.UpdateProviderProvisionStatus)
	}); err != nil {
		return err
	}

	return nil
}

// UpdateProviderProvisionStatus fetches from the REST APIs the latest value for the provider resources
// including vCPU, memory and ephemeral storage, and saves them inside the database.
func (m *Module) UpdateProviderProvisionStatus() error {
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
	statuses, err := m.getProvidersLiveStatuses(addresses)
	if err != nil {
		return err
	}

	return m.db.SaveProviderProvisionStatuses(statuses, height)
}

func (m *Module) getProvidersLiveStatuses(addresses []string) ([]*provider.Status, error) {

	providerStatuses := make([]*provider.Status, len(addresses))
	for i, address := range addresses {
		// Get the provision status of a provider
		status, err := m.source.GetProviderLiveStatus(address)
		if err != nil {
			return nil, err
		}
		providerStatuses[i] = status
	}

	return providerStatuses, nil
}
