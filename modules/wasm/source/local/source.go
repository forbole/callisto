package local

import (
	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/forbole/juno/v3/node/local"

	wasmsource "github.com/forbole/bdjuno/v3/modules/wasm/source"
)

var (
	_ wasmsource.Source = &Source{}
)

// Source implements stakingsource.Source using a local node
type Source struct {
	*local.Source
	q wasmdtypes.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, querier wasmdtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      querier,
	}
}
