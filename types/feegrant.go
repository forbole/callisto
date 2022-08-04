package types

import feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"

// FeeGrant represents the x/feegrant module
type FeeGrant struct {
	feegranttypes.Grant
	Height int64
}

// NewFeeGrant allows to build a new FeeGrant instance
func NewFeeGrant(feegrant feegranttypes.Grant, height int64) FeeGrant {
	return FeeGrant{
		feegrant,
		height,
	}
}

type GrantRemoval struct {
	Grantee string
	Granter string
	Height  int64
}

// NewGrantRemoval allows to build a new GrantRemoval instance
func NewGrantRemoval(grantee string, granter string, height int64) GrantRemoval {
	return GrantRemoval{
		grantee,
		granter,
		height,
	}
}
