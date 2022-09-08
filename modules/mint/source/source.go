package source

import (
	minttypes "github.com/Stride-Labs/stride/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
