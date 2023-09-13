package source

import (
	staketypes "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
)

type Source interface {
	GetStakes(address []string, height int64) ([]*staketypes.Stake, error)
}
