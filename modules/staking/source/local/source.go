package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v2/node/local"

	stakingsource "github.com/forbole/bdjuno/v2/modules/staking/source"
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

// GetDelegation implements stakingsource.Source
func (s Source) GetDelegation(height int64, delegator string, validator string) (stakingtypes.DelegationResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return stakingtypes.DelegationResponse{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Delegation(
		sdk.WrapSDKContext(ctx),
		&stakingtypes.QueryDelegationRequest{DelegatorAddr: delegator, ValidatorAddr: validator},
	)
	if err != nil {
		return stakingtypes.DelegationResponse{}, err
	}

	return *res.DelegationResponse, nil
}

// GetValidatorDelegations implements stakingsource.Source
func (s Source) GetValidatorDelegations(height int64, valOperAddr string) ([]stakingtypes.DelegationResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var delegations []stakingtypes.DelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.ValidatorDelegations(
			sdk.WrapSDKContext(ctx),
			&stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: valOperAddr,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 delegations at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		delegations = append(delegations, res.DelegationResponses...)
	}

	return delegations, nil
}

// GetDelegatorDelegations implements stakingsource.Source
func (s Source) GetDelegatorDelegations(height int64, delegator string) ([]stakingtypes.DelegationResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var delegations []stakingtypes.DelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.DelegatorDelegations(
			sdk.WrapSDKContext(ctx),
			&stakingtypes.QueryDelegatorDelegationsRequest{
				DelegatorAddr: delegator,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 delegations at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		delegations = append(delegations, res.DelegationResponses...)
	}

	return delegations, nil
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

// GetDelegatorRedelegations implements stakingsource.Source
func (s Source) GetDelegatorRedelegations(height int64, delegator string) ([]stakingtypes.RedelegationResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var redelegations []stakingtypes.RedelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.Redelegations(
			sdk.WrapSDKContext(ctx),
			&stakingtypes.QueryRedelegationsRequest{
				DelegatorAddr: delegator,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 delegations at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		redelegations = append(redelegations, res.RedelegationResponses...)
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
func (s Source) GetUnbondingDelegations(height int64, delegator string) ([]stakingtypes.UnbondingDelegation, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var delegations []stakingtypes.UnbondingDelegation
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.DelegatorUnbondingDelegations(
			sdk.WrapSDKContext(ctx),
			&stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
				DelegatorAddr: delegator,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 unbonding delegations at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		delegations = append(delegations, res.UnbondingResponses...)
	}

	return delegations, nil

}
