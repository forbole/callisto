package remote

import (
	"fmt"

	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/forbole/juno/v3/node/remote"

	profilessource "github.com/forbole/bdjuno/v3/modules/profiles/source"
)

var (
	_ profilessource.Source = &Source{}
)

// Source implements profilessource.Source using a remote node
type Source struct {
	*remote.Source
	client profilestypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, client profilestypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// GetParams implements profilessource.Source
func (s *Source) GetParams(height int64) (profilestypes.Params, error) {
	res, err := s.client.Params(remote.GetHeightRequestContext(s.Ctx, height), &profilestypes.QueryParamsRequest{})
	if err != nil {
		return profilestypes.Params{}, fmt.Errorf("error while reading profiles params: %s", err)
	}

	return res.Params, nil
}
