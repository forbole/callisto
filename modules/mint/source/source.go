package source

import (
	minttypes "github.com/MonikaCat/stride/v15/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
