package types

import markertypes "github.com/provenance-io/provenance/x/marker/types"

// MarkerParams represents the x/marker parameters
type MarkerParams struct {
	markertypes.Params
	Height int64
}

// NewMarkerParams allows to build a new MarkerParams instance
func NewMarkerParams(params markertypes.Params, height int64) *MarkerParams {
	return &MarkerParams{
		Params: params,
		Height: height,
	}
}
