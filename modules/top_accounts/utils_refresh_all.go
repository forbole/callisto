package top_accounts

import "fmt"

func (m *Module) RefreshAll(address string) error {
	// Query the latest chain height
	latestHeight, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Set the height 5 blocks lower to avoid error
	// codespace sdk code 26: invalid height: cannot query with height in the future
	height := latestHeight - 5
	
	err = m.bankModule.UpdateBalances([]string{address}, height)
	if err != nil {
		return fmt.Errorf("error while refreshing balance of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshDelegations(address, height)
	if err != nil {
		return fmt.Errorf("error while refreshing delegations of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshUnbondings(address, height)
	if err != nil {
		return fmt.Errorf("error while refreshing unbonding delegations of account %s, error: %s", address, err)
	}

	err = m.distrModule.RefreshDelegatorRewards([]string{address}, height)
	if err != nil {
		return fmt.Errorf("error while refreshing rewards of account %s, error: %s", address, err)
	}

	err = m.refreshTopAccountsSum([]string{address}, height)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum %s, error: %s", address, err)
	}

	return nil
}
