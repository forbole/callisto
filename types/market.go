package types

import markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"

// MarketParams represents the x/market parameters
type MarketParams struct {
	markettypes.Params
	Height int64
}

// NewMarketParams allows to build a new MarketParams instance
func NewMarketParams(params markettypes.Params, height int64) *MarketParams {
	return &MarketParams{
		Params: params,
		Height: height,
	}
}
