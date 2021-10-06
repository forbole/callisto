package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountBalanceHistory contains the information of a given balance at a specific height
type AccountBalanceHistory struct {
	Account       string
	Balance       []sdk.Coin
	Delegations   []sdk.Coin
	Redelegations []sdk.Coin
	Unbonding     []sdk.Coin
	Commission    []sdk.DecCoin
	Reward        []sdk.DecCoin
	Timestamp     time.Time
}

// NewAccountBalanceHistory allows to build a new AccountBalanceHistory instance
func NewAccountBalanceHistory(
	account string, balance, delegations, redelegations, unbonding []sdk.Coin, commission, reward []sdk.DecCoin, timestamp time.Time,
) AccountBalanceHistory {
	return AccountBalanceHistory{
		Account:       account,
		Balance:       balance,
		Delegations:   delegations,
		Redelegations: redelegations,
		Unbonding:     unbonding,
		Commission:    commission,
		Reward:        reward,
		Timestamp:     timestamp,
	}
}
