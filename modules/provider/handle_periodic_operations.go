package provider

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/utils"
	"github.com/forbole/bdjuno/v3/types"

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
	if _, err := scheduler.Every(120).Seconds().Do(func() {
		utils.WatchMethod(m.updateProviderStatus)
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

// updateProviderStatus fetches from the REST APIs the latest value for the provider resources
// including vCPU, memory and ephemeral storage, and saves them inside the database.
func (m *Module) updateProviderStatus() error {
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

	fmt.Println("len(addresses): ", len(addresses))

	invChan := make(chan *types.ProviderStatus)

	goroutineCount := 0
	for _, address := range addresses {
		// Get providers inventory status
		go m.getProviderInventory(address, height, invChan)
		goroutineCount++
		fmt.Println("goroutineCount: ", goroutineCount)
	}

	handledCounter := 0
	i := 0
	ii := 0
	for ch := range invChan {
		if ch.Active {
			fmt.Println(ch)
			i++
			fmt.Println("active count: ", i)
		}

		ii++
		fmt.Println("total count: ", ii)

		handledCounter++
		if handledCounter == len(addresses) {
			close(invChan)
		}
	}

	// return m.db.SaveProviderInventoryStatus(<-invChan, height)
	return nil
}
