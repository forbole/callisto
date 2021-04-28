package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	bstakingutils "github.com/forbole/bdjuno/x/staking/common"

	"github.com/forbole/bdjuno/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
)

// HandleMsg allows to handle the different messages related to the staking module
func HandleMsg(tx *juno.Tx, msg sdk.Msg, cdc codec.Marshaler, db *database.BigDipperDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	// Delegations are handled inside the block handler
	switch cosmosMsg := msg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return handleMsgCreateValidator(tx.Height, cosmosMsg, cdc, db)

	case *stakingtypes.MsgEditValidator:
		return handleEditValidator(tx.Height, cosmosMsg, db)
	}

	return nil
}

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func handleMsgCreateValidator(
	height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(msg.Pubkey, &pubKey)
	if err != nil {
		return err
	}

	operatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	stakingValidator, err := stakingtypes.NewValidator(operatorAddr, pubKey, msg.Description)
	if err != nil {
		return err
	}

	validator, err := bstakingutils.ConvertValidator(cdc, stakingValidator)
	if err != nil {
		return err
	}

	// Save validator
	err = db.SaveValidatorData(validator)
	if err != nil {
		return err
	}

	// Save validator description
	description, err := bstakingutils.GetValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		height,
	)
	if err != nil {
		return err
	}

	err = db.SaveValidatorDescription(description)
	if err != nil {
		return err
	}

	// Save validator commission
	return db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		&msg.Commission.Rate,
		&msg.MinSelfDelegation,
		height,
	))
}

// handleEditValidator handles MsgEditValidator messages, updating the validator info and commission
func handleEditValidator(
	height int64, msg *stakingtypes.MsgEditValidator, db *database.BigDipperDb,
) error {
	// Save validator commission
	err := db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		msg.CommissionRate,
		msg.MinSelfDelegation,
		height,
	))
	if err != nil {
		return err
	}

	// Save validator description
	desc, err := bstakingutils.GetValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		height,
	)
	if err != nil {
		return err
	}

	return db.SaveValidatorDescription(desc)
}
