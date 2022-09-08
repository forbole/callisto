package source

import (
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
