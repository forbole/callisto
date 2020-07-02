package handlers

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	jtypes "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/staking/types"
)

// HandleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func HandleMsgCreateValidator(msg stakingtypes.MsgCreateValidator, db database.BigDipperDb) error {
	stakingValidator := stakingtypes.NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)
	return db.SaveValidatorData(types.NewValidator(
		stakingValidator.GetConsAddr(),
		stakingValidator.GetOperator(),
		stakingValidator.GetConsPubKey(),
		stakingValidator.Description))
}

//HandleEditValidator handles MsgEditValidator
//save the message into the database
func HandleEditValidator(msg stakingtypes.MsgEditValidator, tx jtypes.Tx, db database.BigDipperDb) error {
	commission, err := db.GetCommission(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	if commission.Commission == msg.CommissionRate.Int64() || commission.MinSelfDelegation == msg.MinSelfDelegation.Int64() {
		db.SaveEditCommission(types.NewValidatorCommission(msg.ValidatorAddress,msg.CommissionRate.Int64(),
		msg.MinSelfDelegation.Int64(),tx.Height,tx.TimeStamp))
	}

	db.UpdateValidatorInfo(types.NewValidator(msg.ValidatorAddress.String(), sdk.AccAddress(msg.ValidatorAddress).String(), msg.Description.Moniker,
		msg.Description.Identity, msg.Description.Website,
		msg.Description.SecurityContact, msg.Description.Details))

	return nil
}
