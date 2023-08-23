package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v3/types"
)

type Source interface {
	GetBalances(addresses []string, height int64) ([]types.AccountBalance, error)
	GetSupply(height int64) (sdk.Coins, error)

	// -- For hasura action --
	GetAccountBalance(address string, height int64) ([]sdk.Coin, error)
}
