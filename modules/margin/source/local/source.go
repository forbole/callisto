package local

import (
	"fmt"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marginsource "github.com/forbole/bdjuno/v3/modules/margin/source"
	"github.com/forbole/juno/v3/node/local"
)

var (
	_ marginsource.Source = &Source{}
)

// Source implements marginsource.Source by using a local node
type Source struct {
	*local.Source
	querier margintypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, marginKeeper margintypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: marginKeeper,
	}
}

// GetParams implements marginsource.Source
func (s Source) GetParams(height int64) (*margintypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.GetParams(sdk.WrapSDKContext(ctx), &margintypes.ParamsRequest{})
	if err != nil {
		return nil, nil
	}

	return res.Params, nil
}
