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

	found, err := db.HasValidator(msg.ValidatorAddress.String())
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	found, err = db.HasValidator(msg.DelegatorAddress.String())
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	// Get the delegations
	delegations, err := getDelegations(msg.ValidatorAddress, tx.Height, timestamp, cp)
	if err != nil {
		return err
	}

	// Save the delegations
	return db.SaveCurrentDelegations(delegations)
}

// getDelegations returns the list of all the delegations that the validator having the given address has
// at the given block height (having the given timestamp)
func getDelegations(
	validatorAddress sdk.ValAddress, height int64, timestamp time.Time, cp client.ClientProxy,
) ([]types.Delegation, error) {

	// Handle self delegation
	var responses []staking.DelegationResponse
	endpoint := fmt.Sprintf("/staking/validators/%s/responses?height=%d", validatorAddress.String(), height)
	if _, err := cp.QueryLCDWithHeight(endpoint, &responses); err != nil {
		return nil, err
	}

	delegations := make([]types.Delegation, len(responses))
	for index, delegation := range responses {
		delegations[index] = types.NewDelegation(
			delegation.GetDelegatorAddr(),
			delegation.GetValidatorAddr(),
			delegation.Balance,
			delegation.Shares.String(),
			height,
			timestamp,
		)
	}

	return delegations, nil
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

	return db.SaveHistoricalUnbondingDelegation(types.NewUnbondingDelegation(
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
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	// Store the redelegation
	return db.SaveHistoricalRedelegation(types.NewRedelegation(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		msg.Amount,
		completionTime,
		tx.Height,
		timestamp,
	))
}
