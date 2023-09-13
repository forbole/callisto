package source

import (
	referraltypes "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
)

type Source interface {
	GetStats(dates []string, height int64) ([]*referraltypes.Stats, error)
}
