package source

import (
	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
)

type Source interface {
	Params(height int64) (bundlestypes.Params, error)
}
