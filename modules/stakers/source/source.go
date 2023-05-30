package source

import (
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

type Source interface {
	Params(height int64) (stakerstypes.Params, error)
}
