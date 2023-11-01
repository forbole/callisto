package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	"github.com/forbole/juno/v4/node/local"
	minttypes "github.com/osmosis-labs/osmosis/v20/x/mint/types"
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

// Params implements mintsource.Source
func (s Source) EpochProvisions(height int64) (sdk.Dec, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return sdk.Dec{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.EpochProvisions(sdk.WrapSDKContext(ctx), &minttypes.QueryEpochProvisionsRequest{})
	if err != nil {
		return sdk.Dec{}, nil
	}

	return res.EpochProvisions, nil
}
