package local

import (
	"fmt"

	minttypes "github.com/MonOsmosis/osmosis/v10/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
	"github.com/forbole/juno/v3/node/local"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a local node
type Source struct {
	*local.Source
	querier minttypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier minttypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (minttypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return minttypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &minttypes.QueryParamsRequest{})
	if err != nil {
		return minttypes.Params{}, err
	}

	return res.Params, nil
}
