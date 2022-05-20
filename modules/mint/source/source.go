package source

import (
	minttypes "github.com/osmosis-labs/osmosis/v7/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}
