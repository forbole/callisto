package bank

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
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

	err = m.saveTopAccountsAvailable(balances, height)
	if err != nil {
		return fmt.Errorf("error while saving top accounts available balance: %s", err)
	}

	return nil
}

func (m *Module) saveTopAccountsAvailable(balances []types.NativeTokenBalance, height int64) error {
	if len(balances) == 0 {
		return nil
	}
	err := m.db.SaveTopAccountsBalance("available", balances)
	if err != nil {
		return fmt.Errorf("error while saving top accounts available balances: %s", err)
	}

	var addresses = make([]string, len(balances))
	for index, bal := range balances {
		addresses[index] = bal.Address
	}

	return nil
}
