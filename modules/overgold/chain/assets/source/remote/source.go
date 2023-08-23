package remote

import (
	"github.com/forbole/juno/v3/node/remote"

	"github.com/forbole/bdjuno/v3/modules/overgold/chain/assets/source"

	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client assetstypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, assetClient assetstypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: assetClient,
	}
}

// GetAssets implements Source
func (s Source) GetAssets(addresses []string, height int64) ([]*assetstypes.Asset, error) {
	return []*assetstypes.Asset{}, nil
}
