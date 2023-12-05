package local

import (
	feeexcludertypes "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/forbole/juno/v5/node/local"

	"github.com/forbole/bdjuno/v4/modules/overgold/chain/feeexcluder/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the QueryServer that works on a local node
type Source struct {
	*local.Source
	feeexcluderServer feeexcludertypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, feeexcluderServer feeexcludertypes.QueryServer) *Source {
	return &Source{
		Source:            source,
		feeexcluderServer: feeexcluderServer,
	}
}

// GetFees implements Source
func (s Source) GetFees(denom []string, height int64) ([]*feeexcludertypes.Fees, error) {
	return nil, nil
}
