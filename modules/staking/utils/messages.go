package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"
)

// StoreValidatorFromMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func StoreValidatorFromMsgCreateValidator(
	height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Marshaler, db *database.Db,
) error {
	var pubKey cryptotypes.PubKey
	err := cdc.UnpackAny(msg.Pubkey, &pubKey)
	if err != nil {
		return fmt.Errorf("error while unpacking pub key: %s", err)
	}

	operatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while converting validator address from bech32: %s", err)
	}

	stakingValidator, err := stakingtypes.NewValidator(operatorAddr, pubKey, msg.Description)
	if err != nil {
		return fmt.Errorf("error while creating validator: %s", err)
	}

	validator, err := ConvertValidator(cdc, stakingValidator, height)
	if err != nil {
		return fmt.Errorf("error while converting validator: %s", err)
	}

	desc, err := ConvertValidatorDescription(msg.ValidatorAddress, msg.Description, height)
	if err != nil {
		return fmt.Errorf("error while converting validator description: %s", err)
	}

	// Save the validator
	err = db.SaveValidatorsData([]types.Validator{validator})
	if err != nil {
		return err
	}

	// Save the description
	err = db.SaveValidatorDescription(desc)
	if err != nil {
		return err
	}

	// Save the first self-delegation
	err = db.SaveDelegations([]types.Delegation{
		types.NewDelegation(
			msg.DelegatorAddress,
			msg.ValidatorAddress,
			msg.Value,
			height,
		),
	})
	if err != nil {
		return err
	}

	// Save the commission
	return db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		&msg.Commission.Rate,
		&msg.MinSelfDelegation,
		height,
	))
}

// StoreDelegationFromMessage handles a MsgDelegate and saves the delegation inside the database
func StoreDelegationFromMessage(
	height int64, msg *stakingtypes.MsgDelegate, stakingClient stakingtypes.QueryClient, db *database.Db,
) error {
	header := client.GetHeightRequestHeader(height)
	res, err := stakingClient.Delegation(
		context.Background(),
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: msg.DelegatorAddress,
			ValidatorAddr: msg.ValidatorAddress,
		},
		header,
	)
	if err != nil {
		return fmt.Errorf("error while getting delegation: %s", err)
	}

	delegation := ConvertDelegationResponse(height, *res.DelegationResponse)
	return db.SaveDelegations([]types.Delegation{delegation})
}

// StoreRedelegationFromMessage handles a MsgBeginRedelegate by saving the redelegation inside the database,
// and returns the new redelegation instance
func StoreRedelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate, db *database.Db,
) (*types.Redelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeRedelegate)
	if err != nil {
		return nil, fmt.Errorf("error while searching for EventTypeRedelegate: %s", err)
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, fmt.Errorf("error while searching for AttributeKeyCompletionTime: %s", err)
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, fmt.Errorf("error while covnerting redelegation completion time: %s", err)
	}

	redelegation := types.NewRedelegation(
		msg.DelegatorAddress,
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	return &redelegation, db.SaveRedelegations([]types.Redelegation{redelegation})
}

// StoreUnbondingDelegationFromMessage handles a MsgUndelegate storing the new unbonding delegation inside the database,
// and returns the new unbonding delegation instance
func StoreUnbondingDelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate, db *database.Db,
) (*types.UnbondingDelegation, error) {
	event, err := tx.FindEventByType(index, stakingtypes.EventTypeUnbond)
	if err != nil {
		return nil, fmt.Errorf("error while searching for EventTypeUnbond: %s", err)
	}

	completionTimeStr, err := tx.FindAttributeByKey(event, stakingtypes.AttributeKeyCompletionTime)
	if err != nil {
		return nil, fmt.Errorf("error while searching for AttributeKeyCompletionTime: %s", err)
	}

	completionTime, err := time.Parse(time.RFC3339, completionTimeStr)
	if err != nil {
		return nil, fmt.Errorf("error while covnerting unbonding delegation completion time: %s", err)
	}

	delegation := types.NewUnbondingDelegation(
		msg.DelegatorAddress,
		msg.ValidatorAddress,
		msg.Amount,
		completionTime,
		tx.Height,
	)

	return &delegation, db.SaveUnbondingDelegations([]types.UnbondingDelegation{delegation})
}
