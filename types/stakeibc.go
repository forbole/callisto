package types

import stakeibctypes "github.com/Stride-Labs/stride/x/stakeibc/types"

// StakeIBCParams represents the x/stakeibc parameters
type StakeIBCParams struct {
	stakeibctypes.Params
	Height int64
}

// NewStakeIBCParams allows to build a new StakeIBCParams instance
func NewStakeIBCParams(params stakeibctypes.Params, height int64) *StakeIBCParams {
	return &StakeIBCParams{
		Params: params,
		Height: height,
	}
}
