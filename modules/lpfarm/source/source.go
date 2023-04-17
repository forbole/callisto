package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetTotalLPFarmRewards(farmer string, height int64) (sdk.DecCoins, error)
}
