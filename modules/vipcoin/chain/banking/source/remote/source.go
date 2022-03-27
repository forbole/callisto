/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package remote

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/forbole/juno/v2/node/remote"

	"github.com/forbole/bdjuno/v2/modules/vipcoin/chain/banking/source"
)

var (
	_ source.Source = &Source{}
)

// TODO: other comment
// Source struct
type Source struct {
	*remote.Source
	client banktypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, bankClient banktypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: bankClient,
	}
}

// TODO: change arg and check where is it
// GetBaseTransfers implements Source
func (s Source) GetBaseTransfers(addresses []string, height int64) ([]*bankingtypes.BaseTransfer, error) {
	return []*bankingtypes.BaseTransfer{}, nil
}
