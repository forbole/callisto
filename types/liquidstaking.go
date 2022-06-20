package types

import liquidstakingtypes "github.com/crescent-network/crescent/x/liquidstaking/types"

// LiquidStakingParams represents the x/liquidstaking parameters
type LiquidStakingParams struct {
	liquidstakingtypes.Params
	Height int64
}

// NewLiquidStakingParams allows to build a new LiquidStakingParams instance
func NewLiquidStakingParams(params liquidstakingtypes.Params, height int64) *LiquidStakingParams {
	return &LiquidStakingParams{
		Params: params,
		Height: height,
	}
}
