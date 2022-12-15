package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// AccountBalance represents the balance of an account at a given height
type AccountBalance struct {
	Address string
	Balance sdk.Coins
	Height  int64
}

// NewAccountBalance allows to build a new AccountBalance instance
func NewAccountBalance(address string, balance sdk.Coins, height int64) AccountBalance {
	return AccountBalance{
		Address: address,
		Balance: balance,
		Height:  height,
	}
}

// NativeTokenBalance represents the native token balance of an account at a given height
type NativeTokenBalance struct {
	Address string
	Balance sdk.Int
}

// NewNativeTokenBalance allows to build a new NativeTokenBalance instance
func NewNativeTokenBalance(address string, balance sdk.Int) NativeTokenBalance {
	return NativeTokenBalance{
		Address: address,
		Balance: balance,
	}
}
