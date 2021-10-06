package authority

import (
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v2/modules/utils"
	"github.com/forbole/bdjuno/v2/types"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "authority").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").StartImmediately().Do(func() {
		utils.WatchMethod(m.updateMinGasPrices)
	}); err != nil {
		return err
	}

	return nil
}

// updateMinGasPrices fetches from the REST APIs the latest value for the
// minimum gas prices, and saves it inside the database.
func (m *Module) updateMinGasPrices() error {
	log.Debug().Str("module", "authority").Str("operation", "gas prices").
		Msg("updating minimum gas prices")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the min gas prices
	prices, err := m.source.GetMinimumGasPrices(height)
	if err != nil {
		return err
	}

	return m.db.SaveEMoneyGasPrices(types.NewEMoneyGasPrices(prices, height))
}
