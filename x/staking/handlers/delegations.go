package handlers

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/parse/client"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
)

// HandleMsgDelegate allows to properly handle a MsgDelegate
func HandleMsgDelegate(tx juno.Tx, msg staking.MsgDelegate, db database.BigDipperDb, cp client.ClientProxy) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	if found, _ := db.HasValidator(msg.ValidatorAddress.String()); !found {
		return nil
	}

	if found, _ := db.HasValidator(msg.DelegatorAddress.String()); !found {
		return nil
	}

	if err = saveDelegatorsShares(msg.ValidatorAddress, cp, db, timestamp, tx.Height); err != nil {
		return err
	}

	//for each delegate message it will eventually stored into database
	return db.SaveDelegation(types.NewDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount, tx.Height,
		timestamp,
	))
}

func saveDelegatorsShares(
	validatorAddress sdk.ValAddress, cp client.ClientProxy, db database.BigDipperDb,
	timestamp time.Time, height int64,
) error {

	// Handle self delegation
	var delegations []staking.Delegation
	var delegationsShares []types.DelegationShare

	endpoint := fmt.Sprintf("/staking/validators/%s/delegations?height=%d", validatorAddress.String(), height)
	if _, err := cp.QueryLCDWithHeight(endpoint, &delegations); err != nil {
		return err
	}

	for _, delegation := range delegations {
		delegationsShares = append(delegationsShares, types.NewDelegationShare(
			delegation.GetValidatorAddr(),
			delegation.GetDelegatorAddr(),
			delegation.Shares.Int64(),
			height,
			timestamp,
		))
	}

	return db.SaveDelegationsShares(delegationsShares)
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
