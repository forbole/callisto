package staking

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/forbole/bdjuno/x/utils"

	bstakingutils "github.com/forbole/bdjuno/x/staking/common"

	"github.com/forbole/bdjuno/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
)

// HandleMsg allows to handle the different messages related to the staking module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg, client stakingtypes.QueryClient, cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return handleMsgCreateValidator(tx.Height, cosmosMsg, cdc, db)

	case *stakingtypes.MsgEditValidator:
		return handleEditValidator(tx.Height, cosmosMsg, db)

	case *stakingtypes.MsgDelegate:
		return handleMsgDelegate(tx.Height, cosmosMsg, client, db)

	case *stakingtypes.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(tx, index, cosmosMsg, db)

	case *stakingtypes.MsgUndelegate:
		return handleMsgUndelegate(tx, index, cosmosMsg, db)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func handleMsgCreateValidator(
	height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Marshaler, db *database.BigDipperDb,
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

	validator, err := bstakingutils.ConvertValidator(cdc, stakingValidator)
	if err != nil {
		return err
	}

	// Save validator
	err = db.SaveValidatorData(validator)
	if err != nil {
		return err
	}

	// Save validator description
	description, err := bstakingutils.GetValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		height,
	)
	if err != nil {
		return err
	}

	err = db.SaveValidatorDescription(description)
	if err != nil {
		return err
	}

	// Save validator commission
	return db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		&msg.Commission.Rate,
		&msg.MinSelfDelegation,
		height,
	))
}

// handleEditValidator handles MsgEditValidator messages, updating the validator info and commission
func handleEditValidator(
	height int64, msg *stakingtypes.MsgEditValidator, db *database.BigDipperDb,
) error {
	// Save validator commission
	err := db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		msg.CommissionRate,
		msg.MinSelfDelegation,
		height,
	))
	if err != nil {
		return err
	}

	// Save validator description
	desc, err := bstakingutils.GetValidatorDescription(
		msg.ValidatorAddress,
		msg.Description,
		height,
	)
	if err != nil {
		return err
	}

	return db.SaveValidatorDescription(desc)
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgDelegate handles a MsgDelegate and saves the delegation amount inside the database
func handleMsgDelegate(
	height int64, msg *stakingtypes.MsgDelegate, client stakingtypes.QueryClient, db *database.BigDipperDb,
) error {
	// TODO: Remove the gRPC call when it will be possible to get the new shares amount from the transaction result
	// Cosmos PR: https://github.com/cosmos/cosmos-sdk/pull/9214
	header := utils.GetHeightRequestHeader(height)
	res, err := client.Delegation(
		context.Background(),
		&stakingtypes.QueryDelegationRequest{
			DelegatorAddr: msg.DelegatorAddress,
			ValidatorAddr: msg.ValidatorAddress,
		},
		header,
	)
	if err != nil {
		return err
	}

	delegation := bstakingutils.ConvertDelegationResponse(height, *res.DelegationResponse)
	return db.SaveDelegations([]types.Delegation{delegation})
}

// handleMsgBeginRedelegate handles a MsgBeginRedelegate storing the data inside the database
func handleMsgBeginRedelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate, db *database.BigDipperDb,
) error {
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

// handleMsgUndelegate handles a MsgUndelegate storing the data inside the database
func handleMsgUndelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate, db *database.BigDipperDb,
) error {
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
