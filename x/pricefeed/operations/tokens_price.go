package pricefeed

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/pricefeed/apis"
	"github.com/rs/zerolog/log"
)

// UpdatePrice fetch total amount of coins in the system from RPC and store it into database
func UpdatePrice(db database.BigDipperDb) error {
	log.Debug().
		Str("module", "pricefeed").
		Str("operation", "pricefeed").
		Msg("getting token price and market cap")

	coins, err := apis.GetCoinsList()
	if err != nil {
		return err
	}

	// Get the list of token names to retrieve
	names, err := db.GetTokenNames()
	if err != nil {
		return err
	}

	// Will contain the list of all the coins ids to be fetched
	var ids = make([]string, len(names))

	//find the id of the coins
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

	prices, err := apis.GetTokensPrices(ids)
	if err != nil {
		return err
	}

	return db.SaveTokensPrices(prices, time.Now().UTC())
}
