package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/forbole/juno/v2/node/local"

	distrsource "github.com/forbole/bdjuno/v2/modules/distribution/source"
)

var (
	_ distrsource.Source = &Source{}
)

// Source implements distrsource.Source reading the data from a local node
type Source struct {
	*local.Source
	q distrtypes.QueryServer
}

func NewSource(source *local.Source, keeper distrtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      keeper,
	}
}

// CommunityPool implements distrsource.Source
func (s Source) CommunityPool(height int64) (sdk.DecCoins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.CommunityPool(sdk.WrapSDKContext(ctx), &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return nil, err
	}

	return res.Pool, nil
}

// Params implements distrsource.Source
func (s Source) Params(height int64) (distrtypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return distrtypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Params(sdk.WrapSDKContext(ctx), &distrtypes.QueryParamsRequest{})
	if err != nil {
		return distrtypes.Params{}, err
	}

	return res.Params, nil
}
