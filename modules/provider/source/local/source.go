package local

// import (
// 	"fmt"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	providertypes "github.com/akash-network/akash-api/go/node/provider/v1beta3"
// 	"github.com/cosmos/cosmos-sdk/types/query"
// 	"github.com/forbole/juno/v5/node/local"

// 	providersource "github.com/forbole/bdjuno/v4/modules/provider/source"
// )

// var (
// 	_ providersource.Source = &Source{}
// )

// // Source implements providersource.Source using a local node
// type Source struct {
// 	*local.Source
// 	q providertypes.QueryServer
// }

// // NewSource returns a new Source instance
// func NewSource(source *local.Source, querier providertypes.QueryServer) *Source {
// 	return &Source{
// 		Source: source,
// 		q:      querier,
// 	}
// }

// // GetProvider implements providersource.Source
// func (s Source) GetProvider(height int64, ownerAddress string) (providertypes.Provider, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return providertypes.Provider{}, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.q.Provider(
// 		sdk.WrapSDKContext(ctx),
// 		&providertypes.QueryProviderRequest{
// 			Owner: ownerAddress,
// 		},
// 	)
// 	if err != nil {
// 		return providertypes.Provider{}, fmt.Errorf("error while getting provider: %s", err)
// 	}

// 	return res.Provider, nil
// }

// // GetProviders implements providersource.Source
// func (s Source) GetProviders(height int64) ([]providertypes.Provider, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	var providers []providertypes.Provider
// 	var nextKey []byte
// 	var stop bool
// 	for !stop {
// 		res, err := s.q.Providers(
// 			sdk.WrapSDKContext(ctx),
// 			&providertypes.QueryProvidersRequest{
// 				Pagination: &query.PageRequest{
// 					Key:   nextKey,
// 					Limit: 1000, // Query 1000 providers at a time
// 				},
// 			},
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error while getting providers: %s", err)
// 		}
// 		nextKey = res.Pagination.NextKey
// 		stop = len(res.Pagination.NextKey) == 0
// 		providers = append(providers, res.Providers...)
// 	}

// 	return providers, nil
// }
