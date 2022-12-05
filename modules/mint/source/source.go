package source

import (
	minttypes "github.com/Stride-Labs/stride/v4/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
