/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package source

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
)

// Source - describes an interface for work with banking
type Source interface {
	GetBaseTransfers(addresses []string, height int64) ([]*bankingtypes.BaseTransfer, error)
}
