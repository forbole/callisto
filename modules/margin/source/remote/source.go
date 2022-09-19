package remote

import (
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	marginsource "github.com/forbole/bdjuno/v3/modules/margin/source"
	"github.com/forbole/juno/v3/node/remote"
)

var (
	_ marginsource.Source = &Source{}
)

// Source implements govsource.Source using a remote node
type Source struct {
	*remote.Source
	querier margintypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, querier margintypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetParams implements marginsource.Source
func (s Source) GetParams(height int64) (*margintypes.Params, error) {
	res, err := s.querier.GetParams(remote.GetHeightRequestContext(s.Ctx, height), &margintypes.ParamsRequest{})
	if err != nil {
		return nil, nil
	}

	return res.Params, nil
}
