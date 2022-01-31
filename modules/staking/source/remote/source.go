package remote

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v2/node/remote"

	stakingsource "github.com/forbole/bdjuno/v2/modules/staking/source"
)

var (
	_ stakingsource.Source = &Source{}
)

// Source implements stakingsource.Source using a remote node
type Source struct {
	*remote.Source
	stakingClient stakingtypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, stakingClient stakingtypes.QueryClient) *Source {
	return &Source{
		Source:        source,
		stakingClient: stakingClient,
	}
}

// GetValidator implements stakingsource.Source
func (s Source) GetValidator(height int64, valOper string) (stakingtypes.Validator, error) {
	res, err := s.stakingClient.Validator(
		remote.GetHeightRequestContext(s.Ctx, height),
		&stakingtypes.QueryValidatorRequest{ValidatorAddr: valOper},
	)
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("error while getting validator: %s", err)
	}

	return res.Validator, nil
}

// GetValidatorsWithStatus implements stakingsource.Source
func (s Source) GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var validators []stakingtypes.Validator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.Validators(
			ctx,
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

// GetDelegation implements stakingsource.Source
func (s Source) GetDelegation(height int64, delegator string, valOperAddr string) (stakingtypes.DelegationResponse, error) {
	res, err := s.stakingClient.Delegation(
		remote.GetHeightRequestContext(s.Ctx, height),
		&stakingtypes.QueryDelegationRequest{
			ValidatorAddr: valOperAddr,
			DelegatorAddr: delegator,
		},
	)
	if err != nil {
		return stakingtypes.DelegationResponse{}, err
	}

	return *res.DelegationResponse, nil
}

// GetValidatorDelegations implements stakingsource.Source
func (s Source) GetValidatorDelegations(height int64, validator string) ([]stakingtypes.DelegationResponse, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var delegations []stakingtypes.DelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.ValidatorDelegations(
			ctx,
			&stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: validator,
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
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var delegations []stakingtypes.DelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.DelegatorDelegations(
			ctx,
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

		stop = len(res.Pagination.NextKey) == 0
		delegations = append(delegations, res.DelegationResponses...)
	}

	return delegations, nil
}

// GetDelegatorRedelegations implements stakingsource.Source
func (s Source) GetDelegatorRedelegations(height int64, delegator string) ([]stakingtypes.RedelegationResponse, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var redelegations []stakingtypes.RedelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.Redelegations(
			ctx,
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

		stop = len(res.Pagination.NextKey) == 0
		redelegations = append(redelegations, res.RedelegationResponses...)
	}

	return redelegations, nil
}

// GetPool implements stakingsource.Source
func (s Source) GetPool(height int64) (stakingtypes.Pool, error) {
	res, err := s.stakingClient.Pool(remote.GetHeightRequestContext(s.Ctx, height), &stakingtypes.QueryPoolRequest{})
	if err != nil {
		return stakingtypes.Pool{}, err
	}

	return res.Pool, nil
}

// GetParams implements stakingsource.Source
func (s Source) GetParams(height int64) (stakingtypes.Params, error) {
	res, err := s.stakingClient.Params(remote.GetHeightRequestContext(s.Ctx, height), &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return stakingtypes.Params{}, err
	}

	return res.Params, nil
}

// GetUnbondingDelegations implements stakingsource.Source
func (s Source) GetUnbondingDelegations(height int64, delegator string) ([]stakingtypes.UnbondingDelegation, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var delegations []stakingtypes.UnbondingDelegation
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.DelegatorUnbondingDelegations(
			ctx,
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
