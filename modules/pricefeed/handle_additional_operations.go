package pricefeed

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/types/config"

	historyutils "github.com/forbole/bdjuno/modules/history/utils"
	"github.com/forbole/bdjuno/modules/utils"

	"github.com/forbole/bdjuno/database"

	"github.com/forbole/bdjuno/types"

	"github.com/rs/zerolog/log"
)

// StoreTokens stores the tokens defined inside the given configuration into the database
func StoreTokens(cfg *config.Config, db *database.Db) error {
	log.Debug().Str("module", "pricefeed").Msg("storing tokens")

	var prices []types.TokenPrice
	for _, coin := range cfg.GetPricefeedConfig().GetTokens() {
		// Save the coin as a token with its units
		err := db.SaveToken(coin)
		if err != nil {
			return err
		}

		// Create the price entry
		for _, unit := range coin.Units {
			prices = append(prices, types.NewTokenPrice(unit.Denom, 0, 0, time.Now()))
		}
	}

	err := db.SaveTokensPrices(prices)
	if err != nil {
		return fmt.Errorf("error while storing token prices: %s", err)
	}

	if utils.IsModuleEnabled(cfg, types.HistoryModuleName) {
		return historyutils.UpdatePriceHistory(prices, db)
	}

	return nil
}
