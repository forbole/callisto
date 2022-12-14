package types

import (
	liquidstakingtypes "github.com/crescent-network/crescent/v3/x/liquidstaking/types"
)

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

// LiquidStakingState represents the x/liquidstaking state
type LiquidStakingState struct {
	State  liquidstakingtypes.NetAmountState
	Height int64
}

// NewLiquidStakingState allows to build a new LiquidStakingState instance
func NewLiquidStakingState(state liquidstakingtypes.NetAmountState, height int64) *LiquidStakingState {
	return &LiquidStakingState{
		State:  state,
		Height: height,
	}
}
