package remote

import (
	"fmt"

	providersource "github.com/forbole/bdjuno/v3/modules/provider/source"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"

	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
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

// GetProvider implements providersource.Source
func (s Source) GetProvider(height int64, ownerAddress string) (providertypes.Provider, error) {
	res, err := s.providerClient.Provider(remote.GetHeightRequestContext(s.Ctx, height),
		&providertypes.QueryProviderRequest{
			Owner: ownerAddress,
		},
	)
	if err != nil {
		return providertypes.Provider{}, fmt.Errorf("error while getting provider: %s", err)
	}

	return res.Provider, nil
}

// GetProviders implements providersource.Source
func (s Source) GetProviders(height int64) ([]providertypes.Provider, error) {
	var providers []providertypes.Provider
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.providerClient.Providers(remote.GetHeightRequestContext(s.Ctx, height),
			&providertypes.QueryProvidersRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 1000, // Query 1000 providers at a time
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
