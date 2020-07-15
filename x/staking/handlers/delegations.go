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
	if err=saveDelegatorsShares(validatorAddress,deligatorAddress,cp,db,timestamp,tx.Height);err!=nil{
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


func saveDelegatorsShares(validatorAddress sdk.ValAddress,deligatorAddress sdk.AccAddress,cp client.ClientProxy,db database.BigDipperDb,
	timestamp time.Time,height int64)error{
	//handle self delegation
	var delegations []staking.Delegation
	var delegationstype []types.DelegationShare
	endpoint := fmt.Sprintf("/staking/validators/%s/delegations",validatorAddress.String())
	height, ok := cp.QueryLCDWithHeight(endpoint, &delegations)
	if ok != nil {
		return nil
	}
	for _,delegation := range delegations{
		delegationstype = append(delegationstype,types.NewDelegationShare(
			delegation.GetValidatorAddr(),
			delegation.GetDelegatorAddr(),
			delegation.Shares.Int64(),
			height,
			timestamp))
		}
	if err:=db.SaveDelegationsShares(delegationstype);err!=nil{
						return err
	}
	return nil
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

