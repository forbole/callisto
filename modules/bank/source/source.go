package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetSupply(height int64) (sdk.Coins, error)
}
