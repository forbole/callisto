package bank

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// UpdateBalances updates the balances of the accounts having the given addresses,
// taking the data at the provided height
func (m *Module) UpdateBalances(addresses []string, height int64) error {
	log.Debug().Str("module", "bank").Int64("height", height).Msg("updating balances")

	balances, err := m.keeper.GetBalances(addresses, height)
	if err != nil {
		return fmt.Errorf("error while getting account balances: %s", err)
	}

	err = m.db.SaveTopAccountsBalance("available", balances)
	if err != nil {
		return fmt.Errorf("error while saving top accounts available balances: %s", err)
	}

	return nil
}
