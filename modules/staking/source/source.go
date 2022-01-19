package source

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Source interface {
	GetValidator(height int64, valOper string) (stakingtypes.Validator, error)
	GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, error)
	GetDelegation(height int64, delegator string, validator string) (stakingtypes.DelegationResponse, error)
	GetDelegatorDelegations(height int64, delegator string) ([]stakingtypes.DelegationResponse, error)
	GetDelegatorRedelegations(height int64, delegator string) ([]stakingtypes.RedelegationResponse, error)
	GetPool(height int64) (stakingtypes.Pool, error)
	GetParams(height int64) (stakingtypes.Params, error)
	GetUnbondingDelegations(height int64, delegator string) ([]stakingtypes.UnbondingDelegation, error)
}
