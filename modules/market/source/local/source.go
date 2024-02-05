package local

// import (
// 	"fmt"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	markettypes "github.com/akash-network/akash-api/go/node/market/v1beta4"
// 	"github.com/cosmos/cosmos-sdk/types/query"
// 	"github.com/forbole/juno/v5/node/local"

// 	marketsource "github.com/forbole/bdjuno/v4/modules/market/source"
// )

// var (
// 	_ marketsource.Source = &Source{}
// )

// // Source implements marketsource.Source using a local node
// type Source struct {
// 	*local.Source
// 	q markettypes.QueryServer
// }

// // NewSource returns a new Source instance
// func NewSource(source *local.Source, querier markettypes.QueryServer) *Source {
// 	return &Source{
// 		Source: source,
// 		q:      querier,
// 	}
// }

// // GetActiveLeases implements marketsource.Source
// func (s Source) GetActiveLeases(height int64) ([]markettypes.QueryLeaseResponse, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	var leasesResponse []markettypes.QueryLeaseResponse
// 	var nextKey []byte
// 	var stop bool
// 	for !stop {
// 		res, err := s.q.Leases(sdk.WrapSDKContext(ctx),
// 			&markettypes.QueryLeasesRequest{
// 				Filters: markettypes.LeaseFilters{
// 					// Get only active leases
// 					State: markettypes.LeaseActive.String(),
// 				},
// 				Pagination: &query.PageRequest{
// 					Key:   nextKey,
// 					Limit: 1000, // Query 1000 leases at a time
// 				},
// 			},
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error while getting leases: %s", err)
// 		}
// 		nextKey = res.Pagination.NextKey
// 		stop = len(res.Pagination.NextKey) == 0

// 		leasesResponse = append(leasesResponse, res.Leases...)
// 	}

// 	return leasesResponse, nil
// }
