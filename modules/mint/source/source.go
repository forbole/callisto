package source

import (
	minttypes "github.com/Stride-Labs/stride/v12/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
