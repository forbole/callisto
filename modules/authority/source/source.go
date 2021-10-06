package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetMinimumGasPrices(height int64) (sdk.DecCoins, error)
}
