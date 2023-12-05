package remote

import (
	"github.com/forbole/juno/v5/node/remote"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/core/source"

	coretypes "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client coretypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, coreClient coretypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: coreClient,
	}
}

// GetStats implements Source
func (s Source) GetStats(dates []string, height int64) ([]*coretypes.Stats, error) {
	return []*coretypes.Stats{}, nil
}
