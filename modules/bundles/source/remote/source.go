package remote

import (
	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/forbole/juno/v4/node/remote"

	bundlessource "github.com/forbole/bdjuno/v4/modules/bundles/source"
)

var (
	_ bundlessource.Source = &Source{}
)

// Source implements bundlessource.Source using a remote node
type Source struct {
	*remote.Source
	querier bundlestypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier bundlestypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements bundlessource.Source
func (s Source) Params(height int64) (bundlestypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &bundlestypes.QueryParamsRequest{})
	if err != nil {
		return bundlestypes.Params{}, nil
	}

	return res.Params, nil
}
