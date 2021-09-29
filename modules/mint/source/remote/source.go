package remote

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/desmos-labs/juno/node/remote"

	mintsource "github.com/forbole/bdjuno/modules/mint/source"
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

// GetInflation implements mintsource.Source
func (s Source) GetInflation(height int64) (sdk.Dec, error) {
	res, err := s.querier.Inflation(s.Ctx, &minttypes.QueryInflationRequest{}, remote.GetHeightRequestHeader(height))
	if err != nil {
		return sdk.Dec{}, err
	}

	return res.Inflation, nil
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (minttypes.Params, error) {
	res, err := s.querier.Params(s.Ctx, &minttypes.QueryParamsRequest{}, remote.GetHeightRequestHeader(height))
	if err != nil {
		return minttypes.Params{}, nil
	}

	return res.Params, nil
}
