package history

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v2/modules/pricefeed/coingecko"
	"github.com/forbole/bdjuno/v2/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "history").Msg("setting up periodic tasks")

	// Fetch total supply of tokens every 1hr to store historical price data
	if _, err := scheduler.Every(1).Hour().Do(func() {
		utils.WatchMethod(m.updatePricesHistory)
	}); err != nil {
		return fmt.Errorf("error while setting up history period operations: %s", err)
	}

	return nil
}

// updatePricesHistory fetch total amount of coins in the system from RPC
// and store historical perice data inside the database
func (m *Module) updatePricesHistory() error {
	log.Debug().
		Str("module", "history").
		Str("operation", "history").
		Msg("getting token price and market cap history")

	// Get the list of tokens price id
	ids, err := m.db.GetTokensPriceID()
	if err != nil {
		return fmt.Errorf("error while getting tokens price id: %s", err)
	}

	if len(ids) == 0 {
		log.Debug().Str("module", "history").Msg("no traded tokens price id found")
		return nil
	}

	// Get the tokens prices
	prices, err := coingecko.GetTokensPrices(ids)
	if err != nil {
		return fmt.Errorf("error while getting tokens prices history: %s", err)
	}

	return m.UpdatePricesHistory(prices)

}
