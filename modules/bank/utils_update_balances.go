package bank

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
	"github.com/rs/zerolog/log"
)

// UpdateBalances updates the balances of the accounts having the given addresses,
// taking the data at the provided height
func (m *Module) UpdateBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	log.Debug().Str("module", "bank").Int64("height", height).Msg("updating balances")

	balances, err := m.keeper.GetBalances(addresses, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting account balances: %s", err)
	}

	err = m.db.SaveAccountBalances(balances)
	if err != nil {
		return nil, fmt.Errorf("error while saving account balances: %s", err)
	}

	return balances, nil
}
