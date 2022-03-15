package remote

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v2/node/remote"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"

	providersource "github.com/forbole/bdjuno/v2/modules/provider/source"
)

var (
	_ providersource.Source = &Source{}
)

// Source implements providersource.Source using a remote node
type Source struct {
	*remote.Source
	providerClient providertypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, providerClient providertypes.QueryClient) *Source {
	return &Source{
		Source:         source,
		providerClient: providerClient,
	}
}

// GetProviders implements providersource.Source
func (s Source) GetProviders(height int64) ([]providertypes.Provider, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var providers []providertypes.Provider
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.providerClient.Providers(
			ctx,
			&providertypes.QueryProvidersRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 providers at a time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting providers: %s", err)
		}
		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		providers = append(providers, res.Providers...)
	}

	return providers, nil
}
