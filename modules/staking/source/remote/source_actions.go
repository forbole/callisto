package remote

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v2/node/remote"
)

// GetDelegationsWithPagination implements stakingsource.Source
func (s Source) GetDelegationsWithPagination(
	height int64, delegator string, pagination *query.PageRequest,
) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {

	header := remote.GetHeightRequestHeader(height)

	res, err := s.stakingClient.DelegatorDelegations(
		s.Ctx,
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
			Pagination: &query.PageRequest{
				Limit:      pagination.GetLimit(),
				Offset:     pagination.GetOffset(),
				CountTotal: pagination.GetCountTotal(),
			},
		},
		header,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUnbondingDelegations implements stakingsource.Source
func (s Source) GetUnbondingDelegations(height int64, delegator string, pagination *query.PageRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	header := remote.GetHeightRequestHeader(height)

	unbondingDelegations, err := s.stakingClient.DelegatorUnbondingDelegations(
		s.Ctx,
		&stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: delegator,
			Pagination: &query.PageRequest{
				Limit:      pagination.GetLimit(),
				Offset:     pagination.GetOffset(),
				CountTotal: pagination.GetCountTotal(),
			},
		},
		header,
	)
	if err != nil {
		return nil, err
	}

	return unbondingDelegations, nil
}

// GetRedelegations implements stakingsource.Source
func (s Source) GetRedelegations(height int64, request *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error) {
	header := remote.GetHeightRequestHeader(height)

	redelegations, err := s.stakingClient.Redelegations(
		s.Ctx,
		&stakingtypes.QueryRedelegationsRequest{
			DelegatorAddr: request.DelegatorAddr,
			Pagination: &query.PageRequest{
				Limit:      request.Pagination.Limit,
				Offset:     request.Pagination.Offset,
				CountTotal: request.Pagination.CountTotal,
			},
		},
		header,
	)
	if err != nil {
		return nil, err
	}
	return redelegations, nil
}

// GetValidatorDelegationsWithPagination implements stakingsource.Source
func (s Source) GetValidatorDelegationsWithPagination(
	height int64, validator string, pagination *query.PageRequest,
) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	header := remote.GetHeightRequestHeader(height)

	res, err := s.stakingClient.ValidatorDelegations(
		s.Ctx,
		&stakingtypes.QueryValidatorDelegationsRequest{
			ValidatorAddr: validator,
			Pagination:    pagination,
		},
		header,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
