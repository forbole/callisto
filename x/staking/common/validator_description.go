package common

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/x/staking/keybase"
	"github.com/forbole/bdjuno/x/staking/types"
)

// GetValidatorDescription returns a new types.ValidatorDescription object by fetching the avatar URL
// using the Keybase APIs
func GetValidatorDescription(
	opAddr string, description stakingtypes.Description, height int64,
) (types.ValidatorDescription, error) {
	avatarURL, err := keybase.GetAvatarURL(description.Identity)
	if err != nil {
		return types.ValidatorDescription{}, err
	}

	return types.NewValidatorDescription(opAddr, description, avatarURL, height), nil
}
