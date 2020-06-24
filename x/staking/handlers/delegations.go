package handlers

import (
	"time"

	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
)

// HandleMsgDelegate allows to properly handle a MsgDelegate
func HandleMsgDelegate(tx juno.Tx, msg staking.MsgDelegate, db database.BigDipperDb) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

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
