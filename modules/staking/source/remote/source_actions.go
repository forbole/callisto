package remote

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v5/node/remote"

	"github.com/forbole/bdjuno/v4/utils"
)

// GetDelegationsWithPagination implements stakingsource.Source
func (s Source) GetDelegationsWithPagination(
	height int64, delegator string, pagination *query.PageRequest,
) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	ctx := utils.GetHeightRequestContext(s.Ctx, height)
	res, err := s.stakingClient.DelegatorDelegations(
		ctx,
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
			Pagination:    pagination,
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUnbondingDelegations implements stakingsource.Source
func (s Source) GetUnbondingDelegations(height int64, delegator string, pagination *query.PageRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	ctx := utils.GetHeightRequestContext(s.Ctx, height)

	unbondingDelegations, err := s.stakingClient.DelegatorUnbondingDelegations(
		ctx,
		&stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: delegator,
			Pagination:    pagination,
		},
	)
	if err != nil {
		return nil, err
	}

	return unbondingDelegations, nil
}

// GetRedelegations implements stakingsource.Source
func (s Source) GetRedelegations(height int64, request *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error) {
	ctx := utils.GetHeightRequestContext(s.Ctx, height)

	redelegations, err := s.stakingClient.Redelegations(ctx, request)
	if err != nil {
		return nil, err
	}
	return redelegations, nil
}

// GetValidatorDelegationsWithPagination implements stakingsource.Source
func (s Source) GetValidatorDelegationsWithPagination(
	height int64, validator string, pagination *query.PageRequest,
) (*stakingtypes.QueryValidatorDelegationsResponse, error) {

	res, err := s.stakingClient.ValidatorDelegations(
		remote.GetHeightRequestContext(s.Ctx, height),
		&stakingtypes.QueryValidatorDelegationsRequest{
			ValidatorAddr: validator,
			Pagination:    pagination,
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUnbondingDelegationsFromValidator implements stakingsource.Source
func (s Source) GetUnbondingDelegationsFromValidator(
	height int64, validator string, pagination *query.PageRequest,
) (*stakingtypes.QueryValidatorUnbondingDelegationsResponse, error) {

	unbondingDelegations, err := s.stakingClient.ValidatorUnbondingDelegations(
		remote.GetHeightRequestContext(s.Ctx, height),
		&stakingtypes.QueryValidatorUnbondingDelegationsRequest{
			ValidatorAddr: validator,
			Pagination:    pagination,
		},
	)
	if err != nil {
		return nil, err
	}

	return unbondingDelegations, nil
}
