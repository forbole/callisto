package handlers

import (
	"time"
	"fmt"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/desmos-labs/juno/parse/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMsgDelegate allows to properly handle a MsgDelegate
func HandleMsgDelegate(tx juno.Tx, msg staking.MsgDelegate, db database.BigDipperDb,cp client.ClientProxy) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}
	validatorAddress := msg.ValidatorAddress
	deligatorAddress := msg.DelegatorAddress
	if found, _ := db.HasValidator(validatorAddress.String()); !found {
		return nil
	}
	if found, _ := db.HasValidator(deligatorAddress.String()); !found {
		return nil
	}
	//handle self delegation
	var delegation staking.Delegation
	endpoint := fmt.Sprintf("/staking/delegators/%s/delegations/%s", deligatorAddress.String(), validatorAddress.String())
	height, ok := cp.QueryLCDWithHeight(endpoint, &delegation)
	if ok != nil {
		return nil
	}
	//check if the delegation is self delegation
	selfAddress := sdk.AccAddress(validatorAddress.Bytes())
	if deligatorAddress.Equals(selfAddress) {
		//get current total delegation
		var validator staking.Validator
		endpoint = fmt.Sprintf("/staking/validators/%s", deligatorAddress.String())
		height, ok = cp.QueryLCDWithHeight(endpoint, &validator)
		db.SaveSelfDelegation(types.NewSelfDelegation(msg.ValidatorAddress,delegation.Shares.Int64(),
					float64(delegation.Shares.Int64())/float64(validator.DelegatorShares.Int64())*100,
					height,timestamp))
	}
		//for each delegate message it will eventually stored into database
		return db.SaveDelegation(types.NewDelegation(
			msg.DelegatorAddress,
			msg.ValidatorAddress,
			msg.Amount, tx.Height,
			timestamp,
		))
}

// HandleMsgUndelegate handles properly a MsgUndelegate
func HandleMsgUndelegate(tx juno.Tx, index int, msg staking.MsgUndelegate, db database.BigDipperDb) error {
	// Get completion time
	event, err := tx.FindEventByType(index, staking.EventTypeUnbond)
	if err != nil {
		return err
	}
	completionTimeStr, err := tx.FindAttributeByKey(event, staking.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}
	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	return db.SaveUnbondingDelegation(types.NewUnbondingDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount,
		completionTime,
		tx.Height,
		timestamp,
	))
}

// HandleMsgBeginRedelegate handles properly MsgBeginRedelegate objects
func HandleMsgBeginRedelegate(tx juno.Tx, index int, msg staking.MsgBeginRedelegate, db database.BigDipperDb) error {
	// Get the completion time
	event, err := tx.FindEventByType(index, staking.EventTypeRedelegate)
	if err != nil {
		return err
	}
	completionTimeStr, err := tx.FindAttributeByKey(event, staking.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}
	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	// Build the redelegation object
	reDelegation := types.NewRedelegation(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	// Store the redelegation
	return db.SaveRedelegation(reDelegation)
}
