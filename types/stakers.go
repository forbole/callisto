package types

import stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"

// StakersParams represents the x/stakers parameters
type StakersParams struct {
	stakerstypes.Params
	Height int64
}

// NewStakersParams allows to build a new StakersParams instance
func NewStakersParams(params stakerstypes.Params, height int64) *StakersParams {
	return &StakersParams{
		Params: params,
		Height: height,
	}
}
