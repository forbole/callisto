package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type EmoneyGasPrice struct {
	AuthorityKey string
	GasPrices    sdk.DecCoins
	Height       int64
}

func NewEmoneyGasPrice(authorityKey string, gasPrices sdk.DecCoins, height int64) EmoneyGasPrice {
	return EmoneyGasPrice{
		AuthorityKey: authorityKey,
		GasPrices:    gasPrices,
		Height:       height,
	}
}
