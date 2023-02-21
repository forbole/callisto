package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	creminttypes "github.com/crescent-network/crescent/v4/x/mint/types"

	"github.com/forbole/juno/v4/node/local"

	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a local node
type Source struct {
	*local.Source
	querier creminttypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier creminttypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (creminttypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return creminttypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &creminttypes.QueryParamsRequest{})
	if err != nil {
		return creminttypes.Params{}, err
	}

	return res.Params, nil
}
