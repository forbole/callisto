package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v2/node/local"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"

	providersource "github.com/forbole/bdjuno/v2/modules/provider/source"
)

var (
	_ providersource.Source = &Source{}
)

// Source implements stakingsource.Source using a local node
type Source struct {
	*local.Source
	q providertypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, querier providertypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}

// GetProviders implements providersource.Source
func (s Source) GetProviders(height int64) ([]providertypes.Provider, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var providers []providertypes.Provider
	var nextKey []byte
	var stop bool
	for !stop {
		res, err := s.q.Providers(
			sdk.WrapSDKContext(ctx),
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
