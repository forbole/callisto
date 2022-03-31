package source

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
)

type Source interface {
	GetAccounts(addresses []string, height int64) ([]*accountstypes.Account, error)
}
