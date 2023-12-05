package local

import (
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/juno/v5/node/local"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/bank/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryServer that works on a local node
type Source struct {
	*local.Source
	bankServer types.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bankServer types.QueryServer) *Source {
	return &Source{
		Source:     source,
		bankServer: bankServer,
	}
}
