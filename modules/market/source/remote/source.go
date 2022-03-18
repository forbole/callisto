package remote

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v2/node/remote"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"

	marketsource "github.com/forbole/bdjuno/v2/modules/market/source"
)

var (
	_ marketsource.Source = &Source{}
)

// Source implements providersource.Source using a remote node
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
func (s Source) GetLeases(height int64) ([]markettypes.QueryLeaseResponse, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var responses []markettypes.QueryLeaseResponse
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.marketClient.Leases(
			ctx,
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
