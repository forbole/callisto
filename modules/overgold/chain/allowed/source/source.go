package source

import (
	allowedtypes "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
)

type Source interface {
	GetAddresses(addresses []string, height int64) ([]*allowedtypes.Addresses, error)
}
