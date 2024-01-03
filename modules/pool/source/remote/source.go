package remote

import (
	poolquerytypes "github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/forbole/juno/v5/node/remote"

	poolsource "github.com/forbole/bdjuno/v4/modules/pool/source"
)

var (
	_ poolsource.Source = &Source{}
)

// Source implements poolsource.Source using a remote node
type Source struct {
	*remote.Source
	poolQuerier poolquerytypes.QueryPoolClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, poolQuerier poolquerytypes.QueryPoolClient) *Source {
	return &Source{
		Source:      source,
		poolQuerier: poolQuerier,
	}
}

// Pools implements poolsource.Source
func (s Source) Pools(height int64) ([]poolquerytypes.PoolResponse, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var pools []poolquerytypes.PoolResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.poolQuerier.Pools(
			ctx,
			&poolquerytypes.QueryPoolsRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 pools at time
				},
			},
		)
		if err != nil {
			return nil, err
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		pools = append(pools, res.Pools...)
	}

	return pools, nil
}
