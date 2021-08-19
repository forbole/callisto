package pricefeed

import (
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/pricefeed/coingecko"
	"github.com/forbole/bdjuno/modules/utils"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// CoinGecko to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(scheduler *gocron.Scheduler, db *database.Db) error {
	log.Debug().Str("module", "pricefeed").Msg("setting up periodic tasks")

	// Fetch total supply of token in 30 seconds each
	if _, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updatePrice(db) })
	}); err != nil {
		return err
	}

	return nil
}

// updatePrice fetch total amount of coins in the system from RPC and store it into database
func updatePrice(db *database.Db) error {
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
	names, err := db.GetTokenUnits()
	if err != nil {
		return err
	}

	// Find the id of the coins
	var ids []string
	for _, coin := range coins {
		for _, tradedToken := range names {
			if strings.EqualFold(coin.Symbol, tradedToken) {
				ids = append(ids, coin.ID)
				break
			}
		}
	}

	if len(ids) == 0 {
		log.Debug().Str("module", "pricefeed").Msg("no traded tokens found")
		return nil
	}

	// Get the tokens prices
	prices, err := coingecko.GetTokensPrices(ids)
	if err != nil {
		return err
	}

	// Save the token prices
	return db.SaveTokensPrices(prices)
}
