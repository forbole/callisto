package remote

import (
	"github.com/forbole/juno/v3/node/remote"

	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
	shieldsource "github.com/forbole/bdjuno/v3/modules/shield/source"
)

var (
	_ shieldsource.Source = &Source{}
)

type Source struct {
	*remote.Source
	shieldClient shieldtypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, shieldClient shieldtypes.QueryClient) *Source {
	return &Source{
		Source:       source,
		shieldClient: shieldClient,
	}
}
