package types

import (
	markertypes "github.com/provenance-io/provenance/x/marker/types"
)

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

// Marker represents the x/marker marker account
type Marker struct {
	Address                string
	AccessControl          []markertypes.AccessGrant
	AllowGovernanceControl bool
	Denom                  string
	MarkerType             markertypes.MarkerType
	Status                 markertypes.MarkerStatus
	Supply                 string
	Height                 int64
}

// NewMarker allows to build a new Marker instance
func NewMarker(
	address string,
	accessControl []markertypes.AccessGrant,
	allowGovernanceControl bool,
	denom string,
	markerType markertypes.MarkerType,
	status markertypes.MarkerStatus,
	supply string,
	height int64) *Marker {
	return &Marker{
		Address:                address,
		AccessControl:          accessControl,
		AllowGovernanceControl: allowGovernanceControl,
		Denom:                  denom,
		MarkerType:             markerType,
		Status:                 status,
		Supply:                 supply,
		Height:                 height,
	}
}
