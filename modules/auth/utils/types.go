package utils

import sdk "github.com/cosmos/cosmos-sdk/types"

// --- Types are defined according to genesis file structure ---

type Accounts struct {
	Accounts []AccountDetails `json:"accounts"`
}

type AccountDetails struct {
	AccountType        string             `json:"@type"`
	BaseVestingAccount BaseVestingAccount `json:"base_vesting_account"`
	StartTime          string             `json:"start_time"`
	VestingPeriods     []Period           `json:"vesting_periods"`
}

type BaseVestingAccount struct {
	BaseAccount     BaseAccount `json:"base_account"`
	OriginalVesting []sdk.Coin  `json:"original_vesting"`
	EndTime         string      `json:"end_time"`
}

type BaseAccount struct {
	Address string `json:"address"`
}

type Period struct {
	Length string     `json:"length"`
	Amount []sdk.Coin `json:"amount"`
}
