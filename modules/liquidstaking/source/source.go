package source

import (
	liquidstakingtypes "github.com/crescent-network/crescent/v4/x/liquidstaking/types"
)

type Source interface {
	Params(height int64) (liquidstakingtypes.Params, error)
	GetLiquidValidators(height int64) ([]liquidstakingtypes.LiquidValidatorState, error)
	GetLiquidStakingStates(height int64) (liquidstakingtypes.NetAmountState, error)
}
