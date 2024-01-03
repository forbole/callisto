package remote

import (
	"fmt"

	stakersquerytypes "github.com/KYVENetwork/chain/x/query/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakerssource "github.com/forbole/bdjuno/v4/modules/stakers/source"
	"github.com/forbole/juno/v5/node/remote"
)

var (
	_ stakerssource.Source = &Source{}
)

// Source implements stakerssource.Source using a remote node
type Source struct {
	*remote.Source
	querier        stakerstypes.QueryClient
	stakersQuerier stakersquerytypes.QueryStakersClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier stakerstypes.QueryClient, stakersQuerier stakersquerytypes.QueryStakersClient) *Source {
	return &Source{
		Source:         source,
		querier:        querier,
		stakersQuerier: stakersQuerier,
	}
}

// Params implements stakerssource.Source
func (s Source) Params(height int64) (stakerstypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &stakerstypes.QueryParamsRequest{})
	if err != nil {
		return stakerstypes.Params{}, nil
	}

	return res.Params, nil
}

// Stakers implements stakerssource.Source
func (s Source) Stakers(height int64) ([]stakersquerytypes.FullStaker, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var stakers []stakersquerytypes.FullStaker
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.stakersQuerier.Stakers(
			ctx,
			&stakersquerytypes.QueryStakersRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 stakers at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error when querying stakers %s", err)

		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		stakers = append(stakers, res.Stakers...)
	}
	return stakers, nil
}
