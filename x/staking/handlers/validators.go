package handlers

import (
	"time"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	jtypes "github.com/desmos-labs/juno/types"

)

// HandleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func HandleMsgCreateValidator(msg stakingtypes.MsgCreateValidator, db database.BigDipperDb) error {
	stakingValidator := stakingtypes.NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)
	return db.SaveValidatorData(types.NewValidator(
		stakingValidator.GetConsAddr(),
		stakingValidator.GetOperator(),
		stakingValidator.GetConsPubKey(),
	))
}
func HandleEditValidator(msg stakingtypes.MsgEditValidator, tx jtypes.Tx, db database.BigDipperDb) error {
	address := msg.ValidatorAddress
	if found, _ := db.HasValidator(address.String()); !found {
		return nil
	}

	commission, err := db.GetCommission(msg.ValidatorAddress)
	if err!=nil{
		return err
	}
	if commission.Commission == msg.CommissionRate.Int64() || commission.MinSelfDelegation == msg.MinSelfDelegation.Int64() {
		//change commission table
		commission.Height = tx.Height
		commission.Timestamp,_ = time.Parse(time.RFC3339, tx.Timestamp)
		commission.ValidatorAddress = msg.ValidatorAddress.String()
		db.SaveEditCommission(commission)
	}

	return nil
}
