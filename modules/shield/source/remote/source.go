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

// GetPoolProviders implements shieldsource.Source
func (s Source) GetPoolProviders(height int64) ([]shieldtypes.Provider, error) {
	res, err := s.shieldClient.Providers(
		remote.GetHeightRequestContext(s.Ctx, height),
		&shieldtypes.QueryProvidersRequest{},
	)

	if err != nil {
		return []shieldtypes.Provider{}, err
	}
	return res.Providers, nil
}

// GetShieldStatus implements shieldsource.Source
func (s Source) GetShieldStatus(height int64) (*shieldtypes.QueryShieldStatusResponse, error) {
	res, err := s.shieldClient.ShieldStatus(
		remote.GetHeightRequestContext(s.Ctx, height),
		&shieldtypes.QueryShieldStatusRequest{},
	)

	if err != nil {
		return &shieldtypes.QueryShieldStatusResponse{}, err
	}
	return res, nil
}
