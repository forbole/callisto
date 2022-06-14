package source

import (
	minttypes "github.com/MonOsmosis/osmosis/v10/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
