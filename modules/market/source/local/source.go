package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/local"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"

	marketsource "github.com/forbole/bdjuno/v3/modules/market/source"
)

var (
	_ marketsource.Source = &Source{}
)

// Source implements marketsource.Source using a local node
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
func (s Source) GetLeases(height int64) ([]markettypes.Lease, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var leases []markettypes.Lease
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.q.Leases(sdk.WrapSDKContext(ctx),
			&markettypes.QueryLeasesRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 leases at a time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting leases: %s", err)
		}
		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0

		for _, l := range res.Leases {
			leases = append(leases, l.Lease)
		}
	}

	return leases, nil
}
