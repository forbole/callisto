package remote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/node/remote"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"

	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a remote node
type Source struct {
	*remote.Source
	querier minttypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier minttypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetEpochProvisions implements mintsource.Source
func (s Source) GetEpochProvisions(height int64) (sdk.Dec, error) {
	res, err := s.querier.EpochProvisions(remote.GetHeightRequestContext(s.Ctx, height), &minttypes.QueryEpochProvisionsRequest{})
	if err != nil {
		return sdk.Dec{}, err
	}

	return res.EpochProvisions, nil
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (minttypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &minttypes.QueryParamsRequest{})
	if err != nil {
		return minttypes.Params{}, nil
	}

	return res.Params, nil
}
