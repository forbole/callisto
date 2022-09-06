package types

import margintypes "github.com/Sifchain/sifnode/x/margin/types"

// MarginParams represents the x/margin parameters
type MarginParams struct {
	*margintypes.Params
	Height int64
}

// NewMarginParams allows to build a new MarginParams instance
func NewMarginParams(params *margintypes.Params, height int64) *MarginParams {
	return &MarginParams{
		Params: params,
		Height: height,
	}
}
