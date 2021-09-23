package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type EmoneyGasPrice struct {
	Authority string
	GasPrices sdk.DecCoins
	Height    int64
}

func NewEmoneyGasPrice(authority string, gasPrices sdk.DecCoins, height int64) EmoneyGasPrice {
	return EmoneyGasPrice{
		Authority: authority,
		GasPrices: gasPrices,
		Height:    height,
	}
}
