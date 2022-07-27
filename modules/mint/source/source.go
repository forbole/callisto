package source

import (
	creminttypes "github.com/crescent-network/crescent/v2/x/mint/types"
)

type Source interface {
	Params(height int64) (creminttypes.Params, error)
}
