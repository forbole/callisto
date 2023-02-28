package source

import (
	minttypes "github.com/Stride-Labs/stride/v6/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
