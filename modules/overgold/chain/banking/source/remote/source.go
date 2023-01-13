package remote

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"github.com/forbole/juno/v2/node/remote"

	"github.com/forbole/bdjuno/v2/modules/overgold/chain/banking/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client bankingtypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, bankClient bankingtypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: bankClient,
	}
}

// GetBaseTransfers implements Source
func (s Source) GetBaseTransfers(addresses []string, height int64) ([]*bankingtypes.BaseTransfer, error) {
	return []*bankingtypes.BaseTransfer{}, nil
}
