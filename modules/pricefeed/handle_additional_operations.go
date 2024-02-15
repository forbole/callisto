package pricefeed

import (
	"fmt"
	"time"

	"github.com/forbole/callisto/v4/types"

	"github.com/rs/zerolog/log"
)

// RunAdditionalOperations implements modules.AdditionalOperationsModule
func (m *Module) RunAdditionalOperations() error {
	err := m.checkConfig()
	if err != nil {
		return err
	}

	return m.storeTokens()
}

// checkConfig checks if the module config is valid
func (m *Module) checkConfig() error {
	if m.cfg == nil {
		return fmt.Errorf("pricefeed config is not set but module is enabled")
	}

	return nil
}

// storeTokens stores the tokens defined inside the given configuration into the database
func (m *Module) storeTokens() error {
	log.Debug().Str("module", "pricefeed").Msg("storing tokens")

	var prices []types.TokenPrice
	for _, coin := range m.cfg.Tokens {
		// Save the coin as a token with its units
		err := m.db.SaveToken(coin)
		if err != nil {
			return fmt.Errorf("error while saving token: %s", err)
		}

		// Create the price entry
		for _, unit := range coin.Units {
			// Skip units with empty price ids
			if unit.PriceID == "" {
				continue
			}

			prices = append(prices, types.NewTokenPrice(unit.Denom, 0, 0, time.Time{}))
		}
	}

	err := m.db.SaveTokensPrices(prices)
	if err != nil {
		return fmt.Errorf("error while storing token prices: %s", err)
	}

	return nil
}
