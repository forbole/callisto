package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"

	"github.com/forbole/juno/v4/node/local"

	"github.com/forbole/bdjuno/v3/modules/shield/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q shieldtypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, querier shieldtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}

// GetPoolParams implements shieldsource.Source
func (s Source) GetPoolParams(height int64) (shieldtypes.PoolParams, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return shieldtypes.PoolParams{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.PoolParams(
		sdk.WrapSDKContext(ctx),
		&shieldtypes.QueryPoolParamsRequest{},
	)

	if err != nil {
		return shieldtypes.PoolParams{}, err
	}
	return res.Params, nil
}

// GetPools implements shieldsource.Source
func (s Source) GetPools(height int64) ([]shieldtypes.Pool, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return []shieldtypes.Pool{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Pools(
		sdk.WrapSDKContext(ctx),
		&shieldtypes.QueryPoolsRequest{},
	)

	if err != nil {
		return []shieldtypes.Pool{}, err
	}
	return res.Pools, nil
}

// GetShieldProviders implements shieldsource.Source
func (s Source) GetShieldProviders(height int64) ([]shieldtypes.Provider, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return []shieldtypes.Provider{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Providers(
		sdk.WrapSDKContext(ctx),
		&shieldtypes.QueryProvidersRequest{},
	)

	if err != nil {
		return []shieldtypes.Provider{}, err
	}
	return res.Providers, nil
}

// GetShieldStatus implements shieldsource.Source
func (s Source) GetShieldStatus(height int64) (*shieldtypes.QueryShieldStatusResponse, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return &shieldtypes.QueryShieldStatusResponse{}, fmt.Errorf("error while loading height: %s", err)
	}
	res, err := s.q.ShieldStatus(
		sdk.WrapSDKContext(ctx),
		&shieldtypes.QueryShieldStatusRequest{},
	)

	if err != nil {
		return &shieldtypes.QueryShieldStatusResponse{}, err
	}
	return res, nil
}
