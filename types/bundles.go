package types

import bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"

// BundlesParams represents the x/bundles parameters
type BundlesParams struct {
	bundlestypes.Params
	Height int64
}

// NewBundlesParams allows to build a new BundlesParams instance
func NewBundlesParams(params bundlestypes.Params, height int64) *BundlesParams {
	return &BundlesParams{
		Params: params,
		Height: height,
	}
}
