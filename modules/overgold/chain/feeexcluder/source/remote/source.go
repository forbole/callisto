package remote

import (
	"github.com/forbole/juno/v5/node/remote"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/feeexcluder/source"

	feeexcludertypes "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryClient that works on a remote node
type Source struct {
	*remote.Source
	client feeexcludertypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, feeexcluderClient feeexcludertypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: feeexcluderClient,
	}
}

// GetFees implements Source
func (s Source) GetFees(denom []string, height int64) ([]*feeexcludertypes.Fees, error) {
	return []*feeexcludertypes.Fees{}, nil
}
