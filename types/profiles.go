package types

import (
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

type ProfilesParams struct {
	Nickname profilestypes.NicknameParams
	DTag     profilestypes.DTagParams
	Bio      profilestypes.BioParams
	Oracle   profilestypes.OracleParams
	Height   int64
}

// NewProfilesParams allows to build a new ProfilesParams instance
func NewProfilesParams(params profilestypes.Params, height int64) ProfilesParams {
	return ProfilesParams{
		Nickname: params.Nickname,
		DTag:     params.DTag,
		Bio:      params.Bio,
		Oracle:   params.Oracle,
		Height:   height,
	}
}
