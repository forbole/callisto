package local

import (
	"github.com/forbole/juno/v2/node/local"

	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking/source"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q bankingtypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk bankingtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      bk,
	}
}

// GetBaseTransfers implements keeper.Source
func (s Source) GetBaseTransfers(addresses []string, height int64) ([]*bankingtypes.BaseTransfer, error) {
	return nil, nil
}
