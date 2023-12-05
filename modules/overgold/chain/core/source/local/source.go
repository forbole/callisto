package local

import (
	coretypes "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	"github.com/forbole/juno/v5/node/local"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/core/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryServer that works on a local node
type Source struct {
	*local.Source
	coreServer coretypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, coreServer coretypes.QueryServer) *Source {
	return &Source{
		Source:     source,
		coreServer: coreServer,
	}
}

// GetStats implements Source
func (s Source) GetStats(dates []string, height int64) ([]*coretypes.Stats, error) {
	return nil, nil
}
