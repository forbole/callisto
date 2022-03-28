/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package remote

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"

	vipcoinwalletskeeper "github.com/forbole/bdjuno/v2/modules/vipcoin/chain/wallets/source"
	"github.com/forbole/juno/v2/node/remote"
)

var (
	_ vipcoinwalletskeeper.Source = &Source{}
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
