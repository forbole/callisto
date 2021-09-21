package pricefeed

import (
	"fmt"
	"strings"

	juno "github.com/desmos-labs/juno/types"

	historyutils "github.com/forbole/bdjuno/modules/history/utils"
	"github.com/forbole/bdjuno/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/pricefeed/coingecko"
	"github.com/forbole/bdjuno/modules/utils"
)

// RegisterPeriodicOps returns the AdditionalOperation that periodically runs fetches from
// CoinGecko to make sure that constantly changing data are synced properly.
func RegisterPeriodicOps(cfg juno.Config, scheduler *gocron.Scheduler, db *database.Db) error {
	log.Debug().Str("module", "pricefeed").Msg("setting up periodic tasks")

	// Fetch total supply of token in 30 seconds each
	if _, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
		utils.WatchMethod(func() error { return updatePrice(cfg, db) })
	}); err != nil {
		return fmt.Errorf("error while setting up pricefeed period operations: %s", err)
	}

	return nil
}

// updatePrice fetch total amount of coins in the system from RPC and store it into database
func updatePrice(cfg juno.Config, db *database.Db) error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("getting token price and market cap")

	// Get the list of coins
	coins, err := coingecko.GetCoinsList()
	if err != nil {
		return fmt.Errorf("error while getting coins list: %s", err)
	}

	// Get the list of token units
	units, err := db.GetTokenUnits()
	if err != nil {
		return fmt.Errorf("error while getting token units: %s", err)
	}

	// Find the id of the coins
	var ids []string
	for _, tradedToken := range units {
		// Skip the token if the price id is empty
		if tradedToken.PriceID == "" {
			continue
		}

		for _, coin := range coins {
			if strings.EqualFold(coin.ID, tradedToken.PriceID) {
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
		return fmt.Errorf("error while getting tokens prices: %s", err)
	}

	// Save the token prices
	err = db.SaveTokensPrices(prices)
	if err != nil {
		return fmt.Errorf("error while saving token prices: %s", err)
	}

	if utils.IsModuleEnabled(cfg, types.HistoryModuleName) {
		return historyutils.UpdatePriceHistory(prices, db)
	}

	return nil
}
