package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	marketsource "github.com/forbole/bdjuno/v2/modules/market/source"
	"github.com/forbole/juno/v2/node/local"

	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

var (
	_ marketsource.Source = &Source{}
)

// Source implements providersource.Source using a remote node
type Source struct {
	*local.Source
	q markettypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, querier markettypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}

// GetLeases implements marketsource.Source
func (s Source) GetLeases(height int64) ([]markettypes.QueryLeaseResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var responses []markettypes.QueryLeaseResponse
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.q.Leases(
			sdk.WrapSDKContext(ctx),
			&markettypes.QueryLeasesRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 providers at a time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting leases: %s", err)
		}
		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		responses = append(responses, res.Leases...)
	}

	return responses, nil
}
