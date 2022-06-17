package types

import (
	markertypes "github.com/MonikaCat/provenance/x/marker/types"
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

// MarkerAccount represents the x/marker marker account
type MarkerAccount struct {
	Address                string
	AccessControl          []markertypes.AccessGrant
	AllowGovernanceControl bool
	Denom                  string
	MarkerType             markertypes.MarkerType
	Status                 markertypes.MarkerStatus
	Supply                 []MarkerSupply
	Height                 int64
}

// NewMarkerAccount allows to build a new MarkerAccount instance
func NewMarkerAccount(
	address string,
	accessControl []markertypes.AccessGrant,
	allowGovernanceControl bool,
	denom string,
	markerType markertypes.MarkerType,
	status markertypes.MarkerStatus,
	supply []MarkerSupply,
	height int64) MarkerAccount {
	return MarkerAccount{
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

// MarkerSupply represents the x/marker supply value
type MarkerSupply struct {
	Denom  string
	Amount string
}

// NewMarkerSupply allows to build a new MarkerSupply instance
func NewMarkerSupply(
	denom string,
	amount string) MarkerSupply {
	return MarkerSupply{
		Denom:  denom,
		Amount: amount,
	}
}
