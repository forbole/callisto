package remote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/forbole/juno/v3/node/local"

	profilessource "github.com/forbole/bdjuno/v3/modules/profiles/source"
)

var (
	_ profilessource.Source = &Source{}
)

// Source implements profilessource.Source using a local node
type Source struct {
	*local.Source
	client profilestypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, client profilestypes.QueryServer) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetParams implements profilessource.Source
func (s *Source) GetParams(height int64) (profilestypes.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return profilestypes.Params{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.client.Params(sdk.WrapSDKContext(ctx), &profilestypes.QueryParamsRequest{})
	if err != nil {
		return profilestypes.Params{}, fmt.Errorf("error while reading profiles params: %s", err)
	}

	return res.Params, nil
}
