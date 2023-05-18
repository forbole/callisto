package top_accounts

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (m *Module) refreshTopAccountsSum(addresses []string, height int64) error {
	for _, addr := range addresses {
		sum, err := m.db.GetAccountBalanceSum(addr)
		if err != nil {
			return fmt.Errorf("error while getting account balance sum : %s", err)
		}

		err = m.db.UpdateTopAccountsSum(addr, sum, height)
		if err != nil {
			return fmt.Errorf("error while updating top accounts sum : %s", err)
		}
	}
	return nil
}

func (m *Module) refreshDelegations(delegator string, height int64) func() {
	return func() {
		err := m.stakingModule.RefreshDelegations(delegator, height)
		if err != nil {
			log.Error().Str("module", "top_accounts").Err(err).
				Str("operation", "refresh delegations").Msg("error while refreshing delegations")
		}
	}
}

func (m *Module) refreshUnbondings(delegatorAddr string, height int64) func() {
	return func() {
		err := m.stakingModule.RefreshUnbondings(delegatorAddr, height)
		if err != nil {
			log.Error().Str("module", "top_accounts").Err(err).
				Str("operation", "refresh unbondings").Msg("error while refreshing unbonding delegations")
		}
	}
}

func (m *Module) refreshBalance(address string, height int64) func() {
	return func() {
		err := m.bankModule.UpdateBalances([]string{address}, height)
		if err != nil {
			log.Error().Str("module", "top_accounts").Err(err).
				Str("operation", "update balance").Msg("error while updating account available balances")
		}

		err = m.refreshTopAccountsSum([]string{address}, height)
		if err != nil {
			log.Error().Str("module", "top_accounts").Err(err).
				Str("operation", "update top accounts sum").Msg("error while refreshing top accounts sum while refreshing balance")
		}
	}
}
