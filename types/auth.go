package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// Account represents a chain account
type Account struct {
	Address string
}

// NewAccount builds a new Account instance
func NewAccount(address string) Account {
	return Account{
		Address: address,
	}
}

// --------------- For Vesting Accounts ---------------

// ContinuousVestingAccount represents a Continuous Vesting Account
type ContinuousVestingAccount struct {
	Type            string
	Address         string
	OriginalVesting sdk.Coins
	EndTime         int64
	StartTime       int64
}

// NewContinuousVestingAccount builds a new ContinuousVestingAccount instance
func NewContinuousVestingAccount(account vestingtypes.ContinuousVestingAccount) ContinuousVestingAccount {
	return ContinuousVestingAccount{
		Type:            "ContinuousVestingAccount",
		Address:         account.Address,
		OriginalVesting: account.OriginalVesting,
		EndTime:         account.EndTime,
		StartTime:       account.StartTime,
	}
}

// ContinuousVestingAccount represents a Delayed Vesting Account
type DelayedVestingAccount struct {
	Type            string
	Address         string
	OriginalVesting sdk.Coins
	EndTime         int64
}

// NewDelayedVestingAccount builds a new DelayedVestingAccount instance
func NewDelayedVestingAccount(account vestingtypes.DelayedVestingAccount) DelayedVestingAccount {
	return DelayedVestingAccount{
		Type:            "DelayedVestingAccount",
		Address:         account.Address,
		OriginalVesting: account.OriginalVesting,
		EndTime:         account.EndTime,
	}
}

// PeriodicVestingAccount represents a Periodic Vesting Account
type PeriodicVestingAccount struct {
	Type            string
	Address         string
	OriginalVesting sdk.Coins
	EndTime         int64
	StartTime       int64
	VestingPeriods  []vestingtypes.Period
}

// NewPeriodicVestingAccount builds a new PeriodicVestingAccount instance
func NewPeriodicVestingAccount(account vestingtypes.PeriodicVestingAccount) PeriodicVestingAccount {
	return PeriodicVestingAccount{
		Type:            "PeriodicVestingAccount",
		Address:         account.Address,
		OriginalVesting: account.OriginalVesting,
		EndTime:         account.EndTime,
		StartTime:       account.StartTime,
		VestingPeriods:  account.VestingPeriods,
	}
}
