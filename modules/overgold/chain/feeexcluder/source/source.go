package source

import (
	feeexcludertypes "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
)

type Source interface {
	GetFees(denom []string, height int64) ([]*feeexcludertypes.Fees, error)
}
