package pricefeed

import (
	"fmt"
	"time"

	"github.com/desmos-labs/juno/client"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/pricefeed/coingecko"
	"github.com/forbole/bdjuno/x/utils"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// CoinGecko to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, cp *client.Proxy, db *database.BigDipperDb) error {
	log.Debug().Str("module", "PriceFeed").Msg("setting up periodic tasks")

	// Fetch total supply of token in 30 seconds each
	if _, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updatePrice(db) })
	}); err != nil {
		return err
	}

	return nil
}

// updatePrice fetch total amount of coins in the system from RPC and store it into database
func updatePrice(db *database.BigDipperDb) error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("getting token price and market cap")

	// Get the list of coins
	coins, err := coingecko.GetCoinsList()
	if err != nil {
		return err
	}

	// Get the list of token names to retrieve
	names, err := db.GetTokenNames()
	if err != nil {
		return err
	}

	// Find the id of the coins
	var ids = make([]string, len(names))
	for i := 0; i < len(coins) && len(ids) < len(names); i++ {
		coin := coins[i]

		for _, name := range names {
			if coin.Name == name {
				ids = append(ids, coin.ID)
				break
			}
		}
	}

	if len(ids) == 0 {
		return fmt.Errorf("cannot find tokens from the API: %s", names)
	}

	// Get the tokens prices
	prices, err := coingecko.GetTokensPrices(ids)
	if err != nil {
		return err
	}

	// Save the token prices
	return db.SaveTokensPrices(prices, time.Now().UTC())
}
