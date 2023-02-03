package top_accounts

import "fmt"

func (m *Module) RefreshAll(address string) error {
	err := m.bankModule.UpdateBalances([]string{address}, 0)
	if err != nil {
		return fmt.Errorf("error while refreshing balance of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshDelegations(0, address)
	if err != nil {
		return fmt.Errorf("error while refreshing delegations of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshRedelegations(0, address)
	if err != nil {
		return fmt.Errorf("error while refreshing redelegations of account %s, error: %s", address, err)
	}

	err = m.stakingModule.RefreshUnbondings(0, address)
	if err != nil {
		return fmt.Errorf("error while refreshing unbonding delegations of account %s, error: %s", address, err)
	}

	err = m.distrModule.RefreshDelegatorRewards(0, []string{address})
	if err != nil {
		return fmt.Errorf("error while refreshing rewards of account %s, error: %s", address, err)
	}

	err = m.refreshTopAccountsSum([]string{address})
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum %s, error: %s", address, err)
	}

	return nil
}
