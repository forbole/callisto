package local

import (
	"github.com/forbole/juno/v2/node/local"

	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/accounts/source"

	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q accountstypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk accountstypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      bk,
	}
}

// GetBalances implements keeper.Source
func (s Source) GetAccounts(addresses []string, height int64) ([]*accountstypes.Account, error) {
	return nil, nil
}
