package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v5/node/local"

	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
)

var (
	_ stakingsource.Source = &Source{}
)

// Source implements stakingsource.Source using a local node
type Source struct {
	*local.Source
	q stakingtypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, querier stakingtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}

// GetValidator implements stakingsource.Source
func (s Source) GetValidator(height int64, valOper string) (stakingtypes.Validator, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Validator(sdk.WrapSDKContext(ctx), &stakingtypes.QueryValidatorRequest{ValidatorAddr: valOper})
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("error while reading validator: %s", err)
	}

	return res.Validator, nil
}

// GetDelegationsWithPagination implements stakingsource.Source
func (s Source) GetDelegationsWithPagination(height int64, delegator string, pagination *query.PageRequest) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.DelegatorDelegations(
		sdk.WrapSDKContext(ctx),
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: delegator,
			Pagination: &query.PageRequest{
				Limit:      pagination.GetLimit(),
				Offset:     pagination.GetOffset(),
				CountTotal: pagination.GetCountTotal(),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetRedelegations implements stakingsource.Source
func (s Source) GetRedelegations(height int64, request *stakingtypes.QueryRedelegationsRequest) (*stakingtypes.QueryRedelegationsResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	redelegations, err := s.q.Redelegations(sdk.WrapSDKContext(ctx), request)
	if err != nil {
		return nil, err
	}

	return redelegations, nil
}

// GetValidatorsWithStatus implements stakingsource.Source
func (s Source) GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var validators []stakingtypes.Validator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.Validators(
			sdk.WrapSDKContext(ctx),
			&stakingtypes.QueryValidatorsRequest{
				Status: status,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validators at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validators = append(validators, res.Validators...)
	}

	return validators, nil
}

// GetPool implements stakingsource.Source
func (s Source) GetPool(height int64) (stakingtypes.Pool, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return stakingtypes.Pool{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Pool(
		sdk.WrapSDKContext(ctx),
		&stakingtypes.QueryPoolRequest{},
	)
	if err != nil {
		return stakingtypes.Pool{}, err
	}

	return res.Pool, nil
}

// GetParams implements stakingsource.Source
func (s Source) GetParams(height int64) (stakingtypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return stakingtypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Params(
		sdk.WrapSDKContext(ctx),
		&stakingtypes.QueryParamsRequest{},
	)
	if err != nil {
		return stakingtypes.Params{}, nil
	}

	return res.Params, nil
}

// GetUnbondingDelegations implements stakingsource.Source
func (s Source) GetUnbondingDelegations(height int64, delegator string, pagination *query.PageRequest) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {

	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	unbondingDelegations, err := s.q.DelegatorUnbondingDelegations(
		sdk.WrapSDKContext(ctx),
		&stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: delegator,
			Pagination: &query.PageRequest{
				Limit:      pagination.GetLimit(),
				Offset:     pagination.GetOffset(),
				CountTotal: pagination.GetCountTotal(),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return unbondingDelegations, nil

}

// GetValidatorDelegationsWithPagination implements stakingsource.Source
func (s Source) GetValidatorDelegationsWithPagination(
	height int64, validator string, pagination *query.PageRequest,
) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.ValidatorDelegations(
		sdk.WrapSDKContext(ctx),
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
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	unbondingDelegations, err := s.q.ValidatorUnbondingDelegations(
		sdk.WrapSDKContext(ctx),
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
