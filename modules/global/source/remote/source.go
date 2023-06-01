package remote

import (
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/forbole/juno/v4/node/remote"

	globalsource "github.com/forbole/bdjuno/v4/modules/global/source"
)

var (
	_ globalsource.Source = &Source{}
)

// Source implements globalsource.Source using a remote node
type Source struct {
	*remote.Source
	querier globaltypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier globaltypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// Params implements globalsource.Source
func (s Source) Params(height int64) (globaltypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &globaltypes.QueryParamsRequest{})
	if err != nil {
		return globaltypes.Params{}, nil
	}

	return res.Params, nil
}
