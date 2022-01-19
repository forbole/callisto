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
		s.Ctx,
		&stakingtypes.QueryValidatorRequest{ValidatorAddr: valOper},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return stakingtypes.Validator{}, fmt.Errorf("error while getting validator: %s", err)
	}

	return res.Validator, nil
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
