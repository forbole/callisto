package source

import (
	stakeibctypes "github.com/MonikaCat/stride/v15/x/stakeibc/types"
)

type Source interface {
	Params(height int64) (stakeibctypes.Params, error)
}
