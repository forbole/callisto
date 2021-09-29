package history

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/types"
)

// UpdateAccountBalanceHistory updates the historic balance for the user having the given address
func (m *Module) UpdateAccountBalanceHistory(address string) error {
	block, err := m.db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	return m.UpdateAccountBalanceHistoryWithTime(address, block.Timestamp)
}

// UpdateAccountBalanceHistoryWithTime updates the historic balance for the user having the given address storing it
// associated to the given time
func (m *Module) UpdateAccountBalanceHistoryWithTime(address string, time time.Time) error {
	if !m.cfg.IsModuleEnabled(moduleName) {
		return nil
	}

	// Get the balance
	balance, err := m.db.GetAccountBalance(address)
	if err != nil {
		return fmt.Errorf("error while getting account balance: %s", err)
	}

	delegations, err := m.db.GetUserDelegationsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user delegation amount: %s", err)
	}

	redelegations, err := m.db.GetUserRedelegationsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user redelegation amount: %s", err)
	}

	unbondingDelegations, err := m.db.GetUserUnBondingDelegationsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user unbonding delegation amount: %s", err)
	}

	// Get the distribution data
	rewards, err := m.db.GetUserDelegatorRewardsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user delegation reward: %s", err)
	}

	commission, err := m.db.GetUserValidatorCommissionAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user vlaidator commission amount: %s", err)
	}

	return m.db.SaveAccountBalanceHistory(types.NewAccountBalanceHistory(
		address,
		balance,
		delegations,
		unbondingDelegations,
		redelegations,
		commission,
		rewards,
		time,
	))
}
