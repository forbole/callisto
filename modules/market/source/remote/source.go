package remote

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"

	marketsource "github.com/forbole/bdjuno/v3/modules/market/source"
)

var (
	_ marketsource.Source = &Source{}
)

// Source implements marketsource.Source using a remote node
type Source struct {
	*remote.Source
	marketClient markettypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, marketClient markettypes.QueryClient) *Source {
	return &Source{
		Source:       source,
		marketClient: marketClient,
	}
}

// GetLeases implements marketsource.Source
func (s Source) GetLeases(height int64) ([]markettypes.Lease, error) {
	var leases []markettypes.Lease
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.marketClient.Leases(remote.GetHeightRequestContext(s.Ctx, height),
			&markettypes.QueryLeasesRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 lease at a time
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
