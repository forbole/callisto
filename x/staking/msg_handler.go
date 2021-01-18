package staking

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/forbole/bdjuno/x/staking/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
)

// HandleMsg allows to handle the different messages related to the staking module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg,
	stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return handleMsgCreateValidator(tx, cosmosMsg, cdc, db)

	case *stakingtypes.MsgDelegate:
		return handleMsgDelegate(tx, cosmosMsg, stakingClient, db)

	case *stakingtypes.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(tx, index, cosmosMsg, db)

	case *stakingtypes.MsgUndelegate:
		return handleMsgUndelegate(tx, index, cosmosMsg, db)

	case *stakingtypes.MsgEditValidator:
		return handleEditValidator(tx, cosmosMsg, db)

	}

	return nil
}

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func handleMsgCreateValidator(
	tx *juno.Tx, msg *stakingtypes.MsgCreateValidator, cdc codec.Marshaler, db *database.BigDipperDb,
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

	// Save validator commission
	err = db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		&msg.Commission.Rate,
		&msg.MinSelfDelegation,
		tx.Height,
	))
	if err != nil {
		return err
	}

	// Save validator description
	err = db.SaveValidatorDescription(types.NewValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		tx.Height,
	))
	if err != nil {
		return err
	}

	consAddr, err := stakingValidator.GetConsAddr()
	if err != nil {
		return err
	}

	consPubKey, err := stakingValidator.ConsPubKey()
	if err != nil {
		return err
	}

	// Save validator
	return db.SaveValidatorData(types.NewValidator(
		consAddr.String(),
		stakingValidator.GetOperator().String(),
		consPubKey.String(),
		msg.DelegatorAddress,
		&msg.Commission.MaxChangeRate,
		&msg.Commission.MaxRate,
	))
}

// handleEditValidator handles MsgEditValidator messages, updating the validator info and commission
func handleEditValidator(
	tx *juno.Tx, msg *stakingtypes.MsgEditValidator, db *database.BigDipperDb,
) error {
	// Save validator commission
	err := db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		msg.CommissionRate,
		msg.MinSelfDelegation,
		tx.Height,
	))
	if err != nil {
		return err
	}

	// Save validator description
	return db.SaveValidatorDescription(types.NewValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		tx.Height,
	))
}

// handleMsgDelegate allows to properly handle a MsgDelegate
func handleMsgDelegate(
	tx *juno.Tx, msg *stakingtypes.MsgDelegate, stakingClient stakingtypes.QueryClient, db *database.BigDipperDb,
) error {
	found, err := db.HasValidator(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	found, err = db.HasValidator(msg.DelegatorAddress)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	// Get the delegations
	delegations, err := utils.GetDelegations(msg.ValidatorAddress, tx.Height, stakingClient)
	if err != nil {
		return err
	}

	// Save the delegations
	return db.SaveDelegations(delegations)
}

// handleMsgUndelegate handles properly a MsgUndelegate
func handleMsgUndelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate, db *database.BigDipperDb,
) error {
	// Get completion time
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeUnbond)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	return db.SaveUnbondingDelegations([]types.UnbondingDelegation{
		types.NewUnbondingDelegation(
			msg.DelegatorAddress,
			msg.ValidatorAddress,
			msg.Amount,
			completionTime,
			tx.Height,
		),
	})
}

// handleMsgBeginRedelegate handles properly MsgBeginRedelegate objects
func handleMsgBeginRedelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate, db *database.BigDipperDb,
) error {
	// Get the completion time
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return err
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return err
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return err
	}

	// Store the redelegation
	return db.SaveRedelegations([]types.Redelegation{
		types.NewRedelegation(
			msg.DelegatorAddress,
			msg.ValidatorSrcAddress,
			msg.ValidatorDstAddress,
			msg.Amount,
			completionTime,
			tx.Height,
		),
	})
}
