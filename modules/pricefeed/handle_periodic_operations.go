package pricefeed

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v2/types"

	"github.com/forbole/bdjuno/v2/modules/pricefeed/coingecko"
	"github.com/forbole/bdjuno/v2/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "pricefeed").Msg("setting up periodic tasks")

	// Fetch total supply of tokens every 2 mins
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		utils.WatchMethod(m.updatePrice)
	}); err != nil {
		return fmt.Errorf("error while setting up pricefeed period operations: %s", err)
	}

	// Fetch total supply of tokens every 1hr to store historical price data
	if _, err := scheduler.Every(1).Hour().Do(func() {
		utils.WatchMethod(m.updatePricesHistory)
	}); err != nil {
		return fmt.Errorf("error while setting up history period operations: %s", err)
	}

	return nil
}

// getTokenPrices allows to get the most up-to-date token prices
func (m *Module) getTokenPrices() ([]types.TokenPrice, error) {
	// Get the list of tokens price id
	ids, err := m.db.GetTokensPriceID()
	if err != nil {
		return nil, fmt.Errorf("error while getting tokens price id: %s", err)
	}

	if len(ids) == 0 {
		log.Debug().Str("module", "pricefeed").Msg("no traded tokens price id found")
		return nil, nil
	}

	// Get the tokens prices
	prices, err := coingecko.GetTokensPrices(ids)
	if err != nil {
		return nil, fmt.Errorf("error while getting tokens prices: %s", err)
	}

	return prices, nil
}

// updatePrice fetch total amount of coins in the system from RPC and store it into database
func (m *Module) updatePrice() error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("updating token price and market cap")

	prices, err := m.getTokenPrices()
	if err != nil {
		return err
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
		Msg("updating token price and market cap history")

	prices, err := m.getTokenPrices()
	if err != nil {
		return err
	}

	err = m.db.SaveTokenPricesHistory(prices)
	if err != nil {
		return fmt.Errorf("error while saving token prices history: %s", err)
	}

	return nil
}
