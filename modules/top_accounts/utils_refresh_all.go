package top_accounts

import "fmt"

func (m *Module) RefreshAll(address string) error {
	// Query the latest chain height
	height, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	err = m.bankModule.UpdateBalances([]string{address}, height)
	if err != nil {
		return fmt.Errorf("error while refreshing balance of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshDelegations(height, address)
	if err != nil {
		return fmt.Errorf("error while refreshing delegations of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshRedelegations(height, address)
	if err != nil {
		return fmt.Errorf("error while refreshing redelegations of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshUnbondings(height, address)
	if err != nil {
		return fmt.Errorf("error while refreshing unbonding delegations of account %s, error: %s", address, err)
	}

	err = m.distrModule.RefreshDelegatorRewards(height, []string{address})
	if err != nil {
		return fmt.Errorf("error while refreshing rewards of account %s, error: %s", address, err)
	}

	err = m.refreshTopAccountsSum([]string{address}, height)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum %s, error: %s", address, err)
	}

	return nil
}
