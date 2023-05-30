package source

import (
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
)

type Source interface {
	Params(height int64) (globaltypes.Params, error)
}
