package remote

import (
	stakeibctypes "github.com/MonikaCat/stride/v15/x/stakeibc/types"
	"github.com/forbole/juno/v5/node/remote"

	stakeibcsource "github.com/forbole/bdjuno/v4/modules/stakeibc/source"
)

var (
	_ stakeibcsource.Source = &Source{}
)

// Source implements stakeibcsource.Source using a remote node
type Source struct {
	*remote.Source
	querier stakeibctypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier stakeibctypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements stakeibcsource.Source
func (s Source) Params(height int64) (stakeibctypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &stakeibctypes.QueryParamsRequest{})
	if err != nil {
		return stakeibctypes.Params{}, nil
	}

	return res.Params, nil
}
