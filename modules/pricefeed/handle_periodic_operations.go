package pricefeed

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v2/modules/pricefeed/coingecko"
	"github.com/forbole/bdjuno/v2/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "pricefeed").Msg("setting up periodic tasks")

	// Fetch the token prices every 2 mins
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		utils.WatchMethod(m.updatePrice)
	}); err != nil {
		return fmt.Errorf("error while setting up pricefeed period operations: %s", err)
	}

	// Update the historical token prices every 1 hour
	if _, err := scheduler.Every(1).Hour().Do(func() {
		utils.WatchMethod(m.updatePricesHistory)
	}); err != nil {
		return fmt.Errorf("error while setting up history period operations: %s", err)
	}

	return nil
}

// updatePrice fetch total amount of coins in the system from RPC and store it into database
func (m *Module) updatePrice() error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("getting token price and market cap")

	// Get the list of tokens price id
	ids, err := m.db.GetTokensPriceID()
	if err != nil {
		return fmt.Errorf("error while getting tokens price id: %s", err)
	}

	if len(ids) == 0 {
		log.Debug().Str("module", "pricefeed").Msg("no traded tokens price id found")
		return nil
	}

	// Get the tokens prices
	prices, err := coingecko.GetTokensPrices(ids)
	if err != nil {
		return fmt.Errorf("error while getting tokens prices: %s", err)
	}

	// Save the token prices
	err = m.db.SaveTokensPrices(prices)
	if err != nil {
		return fmt.Errorf("error while saving token prices: %s", err)
	}

	return nil

}

// updatePricesHistory fetches total amount of coins in the system from RPC
// and stores historical perice data inside the database
func (m *Module) updatePricesHistory() error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("getting token price and market cap history")

	// Get the list of tokens price id
	ids, err := m.db.GetTokensPriceID()
	if err != nil {
		return fmt.Errorf("error while getting tokens price id: %s", err)
	}

	if len(ids) == 0 {
		log.Debug().Str("module", "pricefeed").Msg("no traded tokens price id found")
		return nil
	}

	// Get the tokens prices
	prices, err := coingecko.GetTokensPrices(ids)
	if err != nil {
		return fmt.Errorf("error while getting tokens prices history: %s", err)
	}
	for _, price := range prices {
		price.Timestamp = time.Now()
	}

	return m.db.SaveTokenPricesHistory(prices)

}
