package types

import creminttypes "github.com/crescent-network/crescent/v5/x/mint/types"

// MintParams represents the x/mint parameters
type MintParams struct {
	creminttypes.Params
	Height int64
}

// NewMintParams allows to build a new MintParams instance
func NewMintParams(params creminttypes.Params, height int64) *MintParams {
	return &MintParams{
		Params: params,
		Height: height,
	}
}
