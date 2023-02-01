package source

import (
	stakeibctypes "github.com/Stride-Labs/stride/v5/x/stakeibc/types"
)

type Source interface {
	Params(height int64) (stakeibctypes.Params, error)
}
