package auth

import "github.com/forbole/bdjuno/types"

// DB represents a generic database that allows to perform x/auth related operations
type DB interface {
	SaveAccounts(accounts []types.Account) error
}
