package staking

import "github.com/forbole/bdjuno/types"

// DB represents a generic database that allows to perform some x/staking related operations
type DB interface {
	SaveDelegations(delegations []types.Delegation) error
	SaveRedelegations(redelegations []types.Redelegation) error
	SaveUnbondingDelegations(delegations []types.UnbondingDelegation) error
}
