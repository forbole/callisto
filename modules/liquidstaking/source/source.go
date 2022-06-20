package source

import (
	liquidstakingtypes "github.com/crescent-network/crescent/x/liquidstaking/types"
)

type Source interface {
	Params(height int64) (liquidstakingtypes.Params, error)
	GetLiquidValidators(height int64) ([]liquidstakingtypes.LiquidValidatorState, error)
}
