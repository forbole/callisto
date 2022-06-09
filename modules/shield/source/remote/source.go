package remote

import (
	"github.com/forbole/juno/v3/node/remote"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	shieldsource "github.com/forbole/bdjuno/v3/modules/shield/source"
)

var (
	_ shieldsource.Source = &Source{}
)

type Source struct {
	*remote.Source
	shieldClient shieldtypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, shieldClient shieldtypes.QueryClient) *Source {
	return &Source{
		Source:       source,
		shieldClient: shieldClient,
	}
}

// GetPoolParams implements shieldsource.Source
func (s Source) GetPoolParams(height int64) (shieldtypes.PoolParams, error) {
	res, err := s.shieldClient.PoolParams(
		remote.GetHeightRequestContext(s.Ctx, height),
		&shieldtypes.QueryPoolParamsRequest{},
	)

	if err != nil {
		return shieldtypes.PoolParams{}, err
	}
	return res.Params, nil
}

// GetPools implements shieldsource.Source
func (s Source) GetPools(height int64) ([]shieldtypes.Pool, error) {
	res, err := s.shieldClient.Pools(
		remote.GetHeightRequestContext(s.Ctx, height),
		&shieldtypes.QueryPoolsRequest{},
	)

	if err != nil {
		return []shieldtypes.Pool{}, err
	}
	return res.Pools, nil
}
