package handlers

import (
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	jtypes "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
)

// HandleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func HandleMsgCreateValidator(msg stakingtypes.MsgCreateValidator, db database.BigDipperDb) error {
	stakingValidator := stakingtypes.NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)
	return db.SaveSingleValidatorData(types.NewValidator(
		stakingValidator.GetConsAddr(),
		stakingValidator.GetOperator(),
		stakingValidator.GetConsPubKey(),
		stakingValidator.Description,
		sdktypes.AccAddress(stakingValidator.GetConsAddr()),
	))
}

// HandleEditValidator handles MsgEditValidator messages, updating the validator info
func HandleEditValidator(msg stakingtypes.MsgEditValidator, tx jtypes.Tx, db database.BigDipperDb) error {
	validatorinfo, err := db.GetValidatorData(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	if err := db.SaveEditCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		msg.CommissionRate,
		msg.MinSelfDelegation,
		tx.Height,
		timestamp,
	)); err != nil {
		return err
	}

	return db.UpdateValidatorInfo(types.NewValidator(
		validatorinfo.GetConsAddr(),
		validatorinfo.GetOperator(),
		validatorinfo.GetConsPubKey(),
		msg.Description,
		sdktypes.AccAddress(validatorinfo.GetOperator()),
	))
}
