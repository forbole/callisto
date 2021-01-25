package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type AccountBalance struct {
	Address string
	Balance sdk.Coins
	Height  int64
}

func NewAccountBalance(address string, balance sdk.Coins, height int64) AccountBalance {
	return AccountBalance{
		Address: address,
		Balance: balance,
		Height:  height,
	}
}
