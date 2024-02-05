package pricefeed

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/callisto/v4/types"

	"github.com/forbole/callisto/v4/modules/pricefeed/coingecko"
	"github.com/forbole/callisto/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "pricefeed").Msg("setting up periodic tasks")

	// Fetch the token prices every 2 mins
	if _, err := scheduler.Every(2).Minutes().Do(func() {
		utils.WatchMethod(m.UpdatePrice)
	}); err != nil {
		return fmt.Errorf("error while setting up pricefeed period operations: %s", err)
	}

	// Update the historical token prices every 1 hour
	if _, err := scheduler.Every(1).Hour().Do(func() {
		utils.WatchMethod(m.UpdatePricesHistory)
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

// UpdatePrice fetches the total amount of coins in the system from RPC and stores it in database
func (m *Module) UpdatePrice() error {
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

// UpdatePricesHistory fetches total amount of coins in the system from RPC
// and stores historical perice data inside the database
func (m *Module) UpdatePricesHistory() error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("updating token price and market cap history")

	prices, err := m.getTokenPrices()
	if err != nil {
		return err
	}

	// Normally, the last updated value reflects the time when the price was last updated.
	// If price hasn't changed, the returned timestamp will be the same as one hour ago, and it will not
	// be stored in db as it will be a duplicated value.
	// To fix this, we set each price timestamp to be the same as other ones.
	timestamp := time.Now()
	for _, price := range prices {
		price.Timestamp = timestamp
	}

	err = m.db.SaveTokenPricesHistory(prices)
	if err != nil {
		return fmt.Errorf("error while saving token prices history: %s", err)
	}

	return nil
}
