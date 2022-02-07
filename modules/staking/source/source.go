package source

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Source interface {
	GetValidator(height int64, valOper string) (stakingtypes.Validator, error)
	GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, error)
	GetDelegationsWithPagination(height int64, delegator string, pagination *query.PageRequest) (*stakingtypes.QueryDelegatorDelegationsResponse, error)
	GetRedelegations(height int64, request *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error)
	GetPool(height int64) (stakingtypes.Pool, error)
	GetParams(height int64) (stakingtypes.Params, error)
	GetUnbondingDelegations(height int64, delegator string, pagination *query.PageRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error)
	GetValidatorDelegationsWithPagination(height int64, validator string, pagination *query.PageRequest) (*stakingtypes.QueryValidatorDelegationsResponse, error)
	GetUnbondingDelegationsFromValidator(height int64, validator string, pagination *query.PageRequest) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error)
}
