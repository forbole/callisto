package types

import (
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

type ProfilesParams struct {
	NicknameParams profilestypes.NicknameParams
	DTagParams     profilestypes.DTagParams
	BioParams      profilestypes.BioParams
	OracleParams   profilestypes.OracleParams
	Height         int64
}

// NewProfilesParams allows to build a new ProfilesParams instance
func NewProfilesParams(params profilestypes.Params, height int64) *ProfilesParams {
	return &ProfilesParams{
		NicknameParams: params.Nickname,
		DTagParams:     params.DTag,
		BioParams:      params.Bio,
		OracleParams:   params.Oracle,
		Height:         height,
	}
}
