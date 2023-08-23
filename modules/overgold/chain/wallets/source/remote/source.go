package remote

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"

	"github.com/forbole/juno/v3/node/remote"

	overgoldWalletsKeeper "github.com/forbole/bdjuno/v3/modules/overgold/chain/wallets/source"
)

var (
	_ overgoldWalletsKeeper.Source = &Source{}
)

type Source struct {
	*remote.Source
	wallets walletstypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, walletClient walletstypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		wallets: walletClient,
	}
}

// GetWallets implements bankkeeper.Source
func (s Source) GetWallets(addresses []string, height int64) ([]*walletstypes.Wallet, error) {
	return []*walletstypes.Wallet{}, nil
}
