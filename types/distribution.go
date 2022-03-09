package types

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// DistributionParams represents the parameters of the x/distribution module
type DistributionParams struct {
	distrtypes.Params
	Height int64
}

// NewDistributionParams allows to build a new DistributionParams instance
func NewDistributionParams(params distrtypes.Params, height int64) *DistributionParams {
	return &DistributionParams{
		Params: params,
		Height: height,
	}
}
