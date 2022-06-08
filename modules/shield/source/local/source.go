package local

import (
	"fmt"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/juno/v3/node/local"

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