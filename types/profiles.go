package types

import (
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
)

type ProfilesParams struct {
	Params profilestypes.Params
	Height int64
}

// NewProfilesParams allows to build a new ProfilesParams instance
func NewProfilesParams(params profilestypes.Params, height int64) *ProfilesParams {
	return &ProfilesParams{
		Params: params,
		Height: height,
	}
}
