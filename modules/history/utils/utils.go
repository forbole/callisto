package utils

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
)

// UpdateAccountBalanceHistory updates the historic balance for the user having the given address
func UpdateAccountBalanceHistory(address string, db *database.Db) error {
	block, err := db.GetLastBlock()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	return UpdateAccountBalanceHistoryWithTime(address, block.Timestamp, db)
}

// UpdateAccountBalanceHistoryWithTime updates the historic balance for the user having the given address storing it
// associated to the given time
func UpdateAccountBalanceHistoryWithTime(address string, time time.Time, db *database.Db) error {
	// Get the balance
	balance, err := db.GetAccountBalance(address)
	if err != nil {
		return fmt.Errorf("error while getting account balance: %s", err)
	}

	delegations, err := db.GetUserDelegationsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user delegation amount: %s", err)
	}

	redelegations, err := db.GetUserRedelegationsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user redelegation amount: %s", err)
	}

	unbondingDelegations, err := db.GetUserUnBondingDelegationsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user unbonding delegation amount: %s", err)
	}

	// Get the distribution data
	rewards, err := db.GetUserDelegatorRewardsAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user delegation reward: %s", err)
	}

	commission, err := db.GetUserValidatorCommissionAmount(address)
	if err != nil {
		return fmt.Errorf("error while getting user vlaidator commission amount: %s", err)
	}

	return db.SaveAccountBalanceHistory(types.NewAccountBalanceHistory(
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

// UpdatePriceHistory stores the given prices inside the price history table
func UpdatePriceHistory(prices []types.TokenPrice, db *database.Db) error {
	return db.SaveTokenPricesHistory(prices)
}
