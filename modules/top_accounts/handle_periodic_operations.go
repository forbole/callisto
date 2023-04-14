package top_accounts

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "top accounts").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.RefreshTotalAccounts)
	}); err != nil {
		return fmt.Errorf("error while setting up total top accounts periodic operation: %s", err)
	}

	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.RefreshAvailableBalance)
	}); err != nil {
		return fmt.Errorf("error while setting up top accounts periodic operation: %s", err)
	}

	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		utils.WatchMethod(m.RefreshRewards)
	}); err != nil {
		return fmt.Errorf("error while setting up top accounts periodic operation: %s", err)
	}

	return nil
}

// RefreshTotalAccounts refreshes total number of accounts/wallets in database
func (m *Module) RefreshTotalAccounts() error {
	log.Trace().Str("module", "top accounts").Str("operation", "refresh total accounts").
		Msg("refreshing number of all wallets")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	totalAccountsNumber, err := m.authSource.GetTotalNumberOfAccounts(height)
	if err != nil {
		return fmt.Errorf("error while getting total number of accounts: %s", err)
	}

	err = m.db.SaveTotalAccounts(int64(totalAccountsNumber), height)
	if err != nil {
		return fmt.Errorf("error while storing total number of accounts: %s", err)
	}

	return nil
}

// RefreshAvailableBalance refreshes latest available balance in db
func (m *Module) RefreshAvailableBalance() error {
	log.Trace().Str("module", "top accounts").Str("operation", "refresh available balance").
		Msg("refreshing available balance")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	accounts, err := m.authModule.GetAllBaseAccounts(height)
	if err != nil {
		return fmt.Errorf("error while getting base accounts: %s", err)
	}

	if len(accounts) == 0 {
		return nil
	}

	// Store accounts
	err = m.db.SaveAccounts(accounts)
	if err != nil {
		return err
	}

	// Parse addresses to []string
	var addresses []string
	for _, a := range accounts {
		addresses = append(addresses, a.Address)
	}

	err = m.bankModule.UpdateBalances(addresses, height)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts balances, error: %s", err)
	}

	err = m.refreshTopAccountsSum(addresses, height)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum value: %s", err)
	}

	return nil
}

// RefreshRewards refreshes the rewards for all delegators
func (m *Module) RefreshRewards() error {
	log.Trace().Str("module", "top accounts").Str("operation", "refresh rewards").
		Msg("refreshing delegators rewards")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	// Get the delegators
	delegators, err := m.db.GetDelegators()
	if err != nil {
		return fmt.Errorf("error while getting delegators: %s", err)
	}

	if len(delegators) == 0 {
		return nil
	}

	// Refresh rewards
	err = m.distrModule.RefreshDelegatorRewards(delegators, height)
	if err != nil {
		return fmt.Errorf("error while refreshing delegators rewards: %s", err)
	}

	err = m.refreshTopAccountsSum(delegators, height)
	if err != nil {
		return fmt.Errorf("error while refreshing top accounts sum value: %s", err)
	}

	return nil
}
