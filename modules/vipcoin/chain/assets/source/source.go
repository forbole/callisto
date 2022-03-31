package source

import (
	assetstypes "git.ooo.ua/vipcoin/chain/x/assets/types"
)

// Source - describes an interface for work with assets
type Source interface {
	GetAssets(addresses []string, height int64) ([]*assetstypes.Asset, error)
}
