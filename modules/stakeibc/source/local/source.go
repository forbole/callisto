package local

import (
	"fmt"

	stakeibctypes "github.com/Stride-Labs/stride/v11/x/stakeibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/node/local"

	stakeibcsource "github.com/forbole/bdjuno/v4/modules/stakeibc/source"
)

var (
	_ stakeibcsource.Source = &Source{}
)

// Source implements stakeibcsource.Source using a local node
type Source struct {
	*local.Source
	querier stakeibctypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier stakeibctypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements stakeibcsource.Source
func (s Source) Params(height int64) (stakeibctypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return stakeibctypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &stakeibctypes.QueryParamsRequest{})
	if err != nil {
		return stakeibctypes.Params{}, err
	}

	return res.Params, nil
}
