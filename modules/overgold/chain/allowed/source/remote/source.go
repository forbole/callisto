package remote

import (
	"github.com/forbole/juno/v5/node/remote"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/allowed/source"

	allowedtypes "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client allowedtypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, allowedClient allowedtypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: allowedClient,
	}
}

// GetAddresses implements Source
func (s Source) GetAddresses(addresses []string, height int64) ([]*allowedtypes.Addresses, error) {
	return []*allowedtypes.Addresses{}, nil
}
