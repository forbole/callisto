package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type EMoneyGasPrices struct {
	AuthorityKey string
	GasPrices    sdk.DecCoins
	Height       int64
}

func NewEMoneyGasPrices(authorityKey string, gasPrices sdk.DecCoins, height int64) EMoneyGasPrices {
	return EMoneyGasPrices{
		AuthorityKey: authorityKey,
		GasPrices:    gasPrices,
		Height:       height,
	}
}
