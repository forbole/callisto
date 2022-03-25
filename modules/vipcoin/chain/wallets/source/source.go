/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package source

import walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"

type Source interface {
	GetWallets(addresses []string, height int64) ([]*walletstypes.Wallet, error)
}
