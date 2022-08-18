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

// GetActiveLeases implements marketsource.Source
func (s Source) GetActiveLeases(height int64) ([]markettypes.QueryLeaseResponse, error) {
	var leasesResponse []markettypes.QueryLeaseResponse
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.marketClient.Leases(remote.GetHeightRequestContext(s.Ctx, height),
			&markettypes.QueryLeasesRequest{
				Filters: markettypes.LeaseFilters{
					// Get only active leases
					State: markettypes.LeaseActive.String(),
				},
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

		leasesResponse = append(leasesResponse, res.Leases...)

	}

	return leasesResponse, nil
}
