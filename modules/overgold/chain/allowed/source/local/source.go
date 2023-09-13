package local

import (
	allowedtypes "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	"github.com/forbole/juno/v5/node/local"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/allowed/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryServer that works on a local node
type Source struct {
	*local.Source
	allowedServer allowedtypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, allowedServer allowedtypes.QueryServer) *Source {
	return &Source{
		Source:        source,
		allowedServer: allowedServer,
	}
}

// GetAddresses implements Source
func (s Source) GetAddresses(addresses []string, height int64) ([]*allowedtypes.Addresses, error) {
	return nil, nil
}
