package local

import (
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"

	"github.com/forbole/juno/v3/node/local"

	"github.com/forbole/bdjuno/v3/modules/shield/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q shieldtypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, querier shieldtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}
