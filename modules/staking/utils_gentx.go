package staking

import (
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/v4/modules/staking/keybase"
	"github.com/forbole/bdjuno/v4/modules/utils"
	"github.com/forbole/bdjuno/v4/types"
)

// StoreValidatorFromMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func (m *Module) StoreValidatorsFromMsgCreateValidator(height int64, msg *stakingtypes.MsgCreateValidator) error {
	var pubKey cryptotypes.PubKey
	err := m.cdc.UnpackAny(msg.Pubkey, &pubKey)
	if err != nil {
		return fmt.Errorf("error while unpacking pub key: %s", err)
	}
	avatarURL, err := keybase.GetAvatarURL(msg.Description.Identity)
	if err != nil {
		return fmt.Errorf("error while getting Avatar URL: %s", err)
	}

	// Save the validators
	err = m.db.SaveValidatorData(
		types.NewValidator(
			sdk.ConsAddress(pubKey.Address()).String(),
			msg.ValidatorAddress, pubKey.String(),
			msg.DelegatorAddress,
			&msg.Commission.MaxChangeRate,
			&msg.Commission.MaxRate,
			height,
		),
	)
	if err != nil {
		return err
	}

	// For likecoin dual prefix
	opAddr, err := utils.ConvertAddressPrefix("likevaloper", msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while converting to likevaloper prefix: %s", err)
	}

	// Save the descriptions
	err = m.db.SaveValidatorDescription(
		types.NewValidatorDescription(
			opAddr,
			msg.Description,
			avatarURL,
			height,
		),
	)
	if err != nil {
		return err
	}

	// For likecoin dual prefix
	operAddr, err := utils.ConvertAddressPrefix("likevaloper", msg.ValidatorAddress)
	if err != nil {
		return err
	}

	// Save the commissions
	err = m.db.SaveValidatorCommission(
		types.NewValidatorCommission(
			operAddr,
			&msg.Commission.Rate,
			&msg.MinSelfDelegation,
			height,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
