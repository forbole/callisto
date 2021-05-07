package staking

import "github.com/forbole/bdjuno/types"

// DB represents a generic database that allows to perform some x/staking related operations
type DB interface {
	SaveStakingParams(params types.StakingParams) error
	GetStakingParams() (*types.StakingParams, error)
	SaveDelegations(delegations []types.Delegation) error
	SaveRedelegations(redelegations []types.Redelegation) error
	SaveUnbondingDelegations(delegations []types.UnbondingDelegation) error
	SaveValidators(validators []types.Validator) error
}
