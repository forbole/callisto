package source

import (
	coretypes "git.ooo.ua/vipcoin/ovg-chain/x/core/types"
)

type Source interface {
	GetStats(dates []string, height int64) ([]*coretypes.Stats, error)
}
