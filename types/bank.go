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

// NativeTokenAmount represents the native token balance of an account at a given height
type NativeTokenAmount struct {
	Address string
	Balance sdk.Int
	Height  int64
}

// NewNativeTokenAmount allows to build a new NativeTokenAmount instance
func NewNativeTokenAmount(address string, balance sdk.Int, height int64) NativeTokenAmount {
	return NativeTokenAmount{
		Address: address,
		Balance: balance,
		Height:  height,
	}
}
