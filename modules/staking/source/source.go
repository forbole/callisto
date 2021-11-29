package source

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Source interface {
	GetValidator(height int64, valOper string) (stakingtypes.Validator, error)
	GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, error)
	GetDelegation(height int64, delegator string, validator string) (stakingtypes.DelegationResponse, error)
	GetDelegatorDelegations(height int64, delegator string) ([]stakingtypes.DelegationResponse, error)
	GetValidatorDelegations(height int64, validator string) ([]stakingtypes.DelegationResponse, error)
	GetPool(height int64) (stakingtypes.Pool, error)
	GetParams(height int64) (stakingtypes.Params, error)
}
