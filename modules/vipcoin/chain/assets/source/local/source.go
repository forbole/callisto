package local

import (
	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/assets/source"
	"github.com/forbole/juno/v2/node/local"

	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the asset keeper that works on a local node
type Source struct {
	*local.Source
	q assetstypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, ak assetstypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      ak,
	}
}

// GetAssets implements keeper.Source
func (s Source) GetAssets(addresses []string, height int64) ([]*assetstypes.Asset, error) {
	return nil, nil
}
