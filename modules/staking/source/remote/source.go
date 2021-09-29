package remote

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/v2/node/remote"

	stakingsource "github.com/forbole/bdjuno/modules/staking/source"
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

// GetDelegation implements stakingsource.Source
func (s Source) GetDelegation(height int64, delegator string, validator string) (stakingtypes.DelegationResponse, error) {
	res, err := s.stakingClient.Delegation(
		s.Ctx,
		&stakingtypes.QueryDelegationRequest{
			ValidatorAddr: validator,
			DelegatorAddr: delegator,
		},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return stakingtypes.DelegationResponse{}, err
	}

	return *res.DelegationResponse, nil
}

// GetDelegatorDelegations implements stakingsource.Source
func (s Source) GetDelegatorDelegations(height int64, delegator string) ([]stakingtypes.DelegationResponse, error) {
	header := remote.GetHeightRequestHeader(height)

	var delegations []stakingtypes.DelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.DelegatorDelegations(
			s.Ctx,
			&stakingtypes.QueryDelegatorDelegationsRequest{
				DelegatorAddr: delegator,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 delegations at time
				},
			},
			header,
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

// GetValidatorsWithStatus implements stakingsource.Source
func (s Source) GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, error) {
	header := remote.GetHeightRequestHeader(height)

	var validators []stakingtypes.Validator
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakingClient.Validators(
			s.Ctx,
			&stakingtypes.QueryValidatorsRequest{
				Status: status,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validators at time
				},
			},
			header,
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
	res, err := s.stakingClient.Pool(s.Ctx, &stakingtypes.QueryPoolRequest{}, remote.GetHeightRequestHeader(height))
	if err != nil {
		return stakingtypes.Pool{}, err
	}

	return res.Pool, nil
}

// GetParams implements stakingsource.Source
func (s Source) GetParams(height int64) (stakingtypes.Params, error) {
	res, err := s.stakingClient.Params(s.Ctx, &stakingtypes.QueryParamsRequest{}, remote.GetHeightRequestHeader(height))
	if err != nil {
		return stakingtypes.Params{}, err
	}

	return res.Params, nil
}
