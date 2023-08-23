package local

import (
	"github.com/forbole/juno/v3/node/local"

	"github.com/forbole/bdjuno/v3/modules/overgold/chain/wallets/source"

	"git.ooo.ua/vipcoin/chain/x/wallets/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q types.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk types.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      bk,
	}
}

// GetWallets implements keeper.Source
func (s Source) GetWallets(addresses []string, height int64) ([]*types.Wallet, error) {
	return nil, nil
}
