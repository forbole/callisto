package local

// import (
// 	"fmt"

// 	poolquerytypes "github.com/KYVENetwork/chain/x/query/types"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/query"
// 	"github.com/forbole/juno/v4/node/local"

// 	poolsource "github.com/forbole/bdjuno/v4/modules/pool/source"
// )

// var (
// 	_ poolsource.Source = &Source{}
// )

// // Source implements poolsource.Source using a local node
// type Source struct {
// 	*local.Source
// 	poolQuerier poolquerytypes.QueryPoolClient
// }

// // NewSource returns a new Source instace
// func NewSource(source *local.Source, poolQuerier poolquerytypes.QueryPoolClient) *Source {
// 	return &Source{
// 		Source:      source,
// 		poolQuerier: poolQuerier,
// 	}
// }

// // Pools implements poolsource.Source
// func (s Source) Pools(height int64) ([]poolquerytypes.PoolResponse, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	var pools []poolquerytypes.PoolResponse
// 	var nextKey []byte
// 	var stop = false
// 	for !stop {
// 		res, err := s.poolQuerier.Pools(
// 			sdk.WrapSDKContext(ctx),
// 			&poolquerytypes.QueryPoolsRequest{
// 				Pagination: &query.PageRequest{
// 					Key:   nextKey,
// 					Limit: 100, // Query 100 pools at time
// 				},
// 			},
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		nextKey = res.Pagination.NextKey
// 		stop = len(res.Pagination.NextKey) == 0
// 		pools = append(pools, res.Pools...)
// 	}

// 	return pools, nil
// }
