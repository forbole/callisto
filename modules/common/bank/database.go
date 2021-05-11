package bank

import "github.com/forbole/bdjuno/types"

// DB represents a generic database that supports specific x/bank operations
type DB interface {
	SaveAccountBalances(balances []types.AccountBalance) error
}
