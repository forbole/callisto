package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type EMoneyGasPrices struct {
	GasPrices sdk.DecCoins
	Height    int64
}

func NewEMoneyGasPrices(gasPrices sdk.DecCoins, height int64) EMoneyGasPrices {
	return EMoneyGasPrices{
		GasPrices: gasPrices,
		Height:    height,
	}
}
