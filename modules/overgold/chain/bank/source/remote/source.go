package remote

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/juno/v5/node/remote"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/bank/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client banktypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, bankClient banktypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: bankClient,
	}
}
