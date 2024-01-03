package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/forbole/juno/v5/node/local"

	globalsource "github.com/forbole/bdjuno/v4/modules/global/source"
)

var (
	_ globalsource.Source = &Source{}
)

// Source implements globalsource.Source using a local node
type Source struct {
	*local.Source
	querier globaltypes.QueryServer
}

// NewSource returns a new Source instace
func NewSource(source *local.Source, querier globaltypes.QueryServer) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements globalsource.Source
func (s Source) Params(height int64) (globaltypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return globaltypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.querier.Params(sdk.WrapSDKContext(ctx), &globaltypes.QueryParamsRequest{})
	if err != nil {
		return globaltypes.Params{}, err
	}

	return res.Params, nil
}
