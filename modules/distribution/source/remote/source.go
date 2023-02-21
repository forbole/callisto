package remote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/forbole/juno/v4/node/remote"

	distrsource "github.com/forbole/bdjuno/v3/modules/distribution/source"
)

var (
	_ distrsource.Source = &Source{}
)

// Source implements distrsource.Source querying the data from a remote node
type Source struct {
	*remote.Source
	distrClient distrtypes.QueryClient
}

// NewSource returns a new Source instace
func NewSource(source *remote.Source, distrClient distrtypes.QueryClient) *Source {
	return &Source{
		Source:      source,
		distrClient: distrClient,
	}
}

// CommunityPool implements distrsource.Source
func (s Source) CommunityPool(height int64) (sdk.DecCoins, error) {
	res, err := s.distrClient.CommunityPool(
		remote.GetHeightRequestContext(s.Ctx, height),
		&distrtypes.QueryCommunityPoolRequest{},
	)
	if err != nil {
		return nil, err
	}

	return res.Pool, nil
}

// Params implements distrsource.Source
func (s Source) Params(height int64) (distrtypes.Params, error) {
	res, err := s.distrClient.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&distrtypes.QueryParamsRequest{},
	)
	if err != nil {
		return distrtypes.Params{}, err
	}

	return res.Params, nil
}
