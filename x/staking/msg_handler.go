package staking

import (
	"time"

	cosmosstakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/client"

	stakingtypes "github.com/forbole/bdjuno/x/staking/types"
	"github.com/forbole/bdjuno/x/staking/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
)

// HandleMsg allows to handle the different messages related to the staking module
func HandleMsg(tx types.Tx, index int, msg sdk.Msg, cp *client.Proxy, db *database.BigDipperDb) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "staking").
			Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case staking.MsgCreateValidator:
		return handleMsgCreateValidator(tx, cosmosMsg, db)

	case staking.MsgDelegate:
		return handleMsgDelegate(tx, cosmosMsg, cp, db)

	case staking.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(tx, index, cosmosMsg, db)

	case staking.MsgUndelegate:
		return handleMsgUndelegate(tx, index, cosmosMsg, db)

	case staking.MsgEditValidator:
		return handleEditValidator(tx, cosmosMsg, db)

	}

	return nil
}

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func handleMsgCreateValidator(tx types.Tx, msg cosmosstakingtypes.MsgCreateValidator, db *database.BigDipperDb) error {
	stakingValidator := cosmosstakingtypes.NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	// Save validator commission
	err = db.SaveValidatorCommission(stakingtypes.NewValidatorCommission(
		msg.ValidatorAddress,
		&msg.Commission.Rate,
		&msg.MinSelfDelegation,
		tx.Height,
		timestamp,
	))
	if err != nil {
		return err
	}

	// Save validator description
	err = db.SaveValidatorDescription(stakingtypes.NewValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		tx.Height,
		timestamp,
	))
	if err != nil {
		return err
	}

	// Save validator
	return db.SaveValidatorData(stakingtypes.NewValidator(
		stakingValidator.GetConsAddr(),
		stakingValidator.GetOperator(),
		stakingValidator.GetConsPubKey(),
		sdk.AccAddress(stakingValidator.GetConsAddr()),
		&msg.Commission.MaxChangeRate,
		&msg.Commission.MaxRate,
	))
}

// handleEditValidator handles MsgEditValidator messages, updating the validator info and commission
func handleEditValidator(tx types.Tx, msg cosmosstakingtypes.MsgEditValidator, db *database.BigDipperDb) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	// Save validator commission
	err = db.SaveValidatorCommission(stakingtypes.NewValidatorCommission(
		msg.ValidatorAddress,
		msg.CommissionRate,
		msg.MinSelfDelegation,
		tx.Height,
		timestamp,
	))
	if err != nil {
		return err
	}

	// Save validator description
	return db.SaveValidatorDescription(stakingtypes.NewValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		tx.Height,
		timestamp,
	))
}

// handleMsgDelegate allows to properly handle a MsgDelegate
func handleMsgDelegate(tx types.Tx, msg staking.MsgDelegate, cp *client.Proxy, db *database.BigDipperDb) error {
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
	delegations, err := utils.GetDelegations(msg.ValidatorAddress, tx.Height, timestamp, cp)
	if err != nil {
		return err
	}

	// Save the delegations
	return db.SaveCurrentDelegations(delegations)
}

// handleMsgUndelegate handles properly a MsgUndelegate
func handleMsgUndelegate(tx types.Tx, index int, msg staking.MsgUndelegate, db *database.BigDipperDb) error {
	// Get completion time
	event, err := tx.FindEventByType(index, cosmosstakingtypes.EventTypeUnbond)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, cosmosstakingtypes.AttributeKeyCompletionTime)
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

	return db.SaveHistoricalUnbondingDelegation(stakingtypes.NewUnbondingDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount,
		completionTime,
		tx.Height,
		timestamp,
	))
}

// handleMsgBeginRedelegate handles properly MsgBeginRedelegate objects
func handleMsgBeginRedelegate(tx types.Tx, index int, msg staking.MsgBeginRedelegate, db *database.BigDipperDb) error {
	// Get the completion time
	event, err := tx.FindEventByType(index, cosmosstakingtypes.EventTypeRedelegate)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, cosmosstakingtypes.AttributeKeyCompletionTime)
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
	return db.SaveHistoricalRedelegation(stakingtypes.NewRedelegation(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		msg.Amount,
		completionTime,
		tx.Height,
		timestamp,
	))
}
