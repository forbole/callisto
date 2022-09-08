package types

import icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

// InterchainStakingParams represents the x/interchainstaking parameters
type InterchainStakingParams struct {
	icstypes.Params
	Height int64
}

// NewInterchainStakingParams allows to build a new InterchainStakingParams instance
func NewInterchainStakingParams(params icstypes.Params, height int64) *InterchainStakingParams {
	return &InterchainStakingParams{
		Params: params,
		Height: height,
	}
}
