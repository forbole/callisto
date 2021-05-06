package common

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/modules/bigdipper/staking/keybase"
	bstakingtypes "github.com/forbole/bdjuno/modules/bigdipper/staking/types"
)

// GetValidatorDescription returns a new types.ValidatorDescription object by fetching the avatar URL
// using the Keybase APIs
func GetValidatorDescription(
	opAddr string, description stakingtypes.Description, height int64,
) (bstakingtypes.ValidatorDescription, error) {
	avatarURL, err := keybase.GetAvatarURL(description.Identity)
	if err != nil {
		return bstakingtypes.ValidatorDescription{}, err
	}

	return bstakingtypes.NewValidatorDescription(opAddr, description, avatarURL, height), nil
}
