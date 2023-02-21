package remote

import (
	creminttypes "github.com/crescent-network/crescent/v4/x/mint/types"
	"github.com/forbole/juno/v4/node/remote"

	mintsource "github.com/forbole/bdjuno/v3/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a remote node
type Source struct {
	*remote.Source
	querier creminttypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier creminttypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (creminttypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &creminttypes.QueryParamsRequest{})
	if err != nil {
		return creminttypes.Params{}, nil
	}

	return res.Params, nil
}
