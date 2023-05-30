package types

import globaltypes "github.com/KYVENetwork/chain/x/global/types"

// GlobalParams represents the x/global parameters
type GlobalParams struct {
	globaltypes.Params
	Height int64
}

// NewGlobalParams allows to build a new GlobalParams instance
func NewGlobalParams(params globaltypes.Params, height int64) *GlobalParams {
	return &GlobalParams{
		Params: params,
		Height: height,
	}
}
