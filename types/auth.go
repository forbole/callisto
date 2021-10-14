package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

// --------------- For Vesting Account ---------------

// Account represents a chain account
type VestingAccount struct {
	Address         string
	OriginalVesting sdk.Coins
	EndTime         string
	StartTime       string
	VestingPeriods  []VestingPeriod
}

// NewAccount builds a new VestingAccount instance
func NewVestingAccount(address string, originalVesting sdk.Coins, endTime string, startTime string, vestingPeriods []VestingPeriod) VestingAccount {
	return VestingAccount{
		Address:         address,
		OriginalVesting: originalVesting,
		EndTime:         endTime,
		StartTime:       startTime,
		VestingPeriods:  vestingPeriods,
	}
}

func (u VestingAccount) Equal(v VestingAccount) bool {
	for index, periodU := range u.VestingPeriods {
		periodV := v.VestingPeriods[index]
		if periodU.Amounts.String() != periodV.Amounts.String() {
			return false
		}
	}
	return u.Address == v.Address &&
		u.OriginalVesting.String() == v.OriginalVesting.String() &&
		u.EndTime == v.EndTime &&
		u.StartTime == v.StartTime
}

type VestingPeriod struct {
	Length  string
	Amounts sdk.Coins
}

func NewVestingPeriod(length string, amounts sdk.Coins) VestingPeriod {
	return VestingPeriod{
		Length:  length,
		Amounts: amounts,
	}
}
