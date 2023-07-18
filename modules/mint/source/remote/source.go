package remote

import (
	minttypes "github.com/Stride-Labs/stride/v12/x/mint/types"
	"github.com/forbole/juno/v5/node/remote"

	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
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

// Params implements mintsource.Source
func (s Source) Params(height int64) (minttypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &minttypes.QueryParamsRequest{})
	if err != nil {
		return minttypes.Params{}, nil
	}

	return res.Params, nil
}
