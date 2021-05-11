package staking

import (
	"time"

	"github.com/forbole/bdjuno/modules/common/staking"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	bbankcommon "github.com/forbole/bdjuno/modules/bigdipper/bank/common"
	bstakingcommon "github.com/forbole/bdjuno/modules/bigdipper/staking/common"
	bstakingtypes "github.com/forbole/bdjuno/modules/bigdipper/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to handle the different utils related to the staking module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg, stakingClient stakingtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *bigdipperdb.Db,
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
		return staking.StoreDelegationFromMessage(tx.Height, cosmosMsg, stakingClient, db)

	case *stakingtypes.MsgBeginRedelegate:
		return handleMsgBeginRedelegate(tx, index, cosmosMsg, stakingClient, db)

	case *stakingtypes.MsgUndelegate:
		return handleMsgUndelegate(tx, index, cosmosMsg, stakingClient, bankClient, db)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func handleMsgCreateValidator(
	height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	err := staking.StoreValidatorFromMsgCreateValidator(height, msg, cdc, db)
	if err != nil {
		return err
	}

	// Save validator description
	description, err := bstakingcommon.GetValidatorDescription(
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
	return db.SaveValidatorCommission(bstakingtypes.NewValidatorCommission(
		msg.ValidatorAddress,
		&msg.Commission.Rate,
		&msg.MinSelfDelegation,
		height,
	))
}

// handleEditValidator handles MsgEditValidator utils, updating the validator info and commission
func handleEditValidator(
	height int64, msg *stakingtypes.MsgEditValidator, db *bigdipperdb.Db,
) error {
	// Save validator commission
	err := db.SaveValidatorCommission(bstakingtypes.NewValidatorCommission(
		msg.ValidatorAddress,
		msg.CommissionRate,
		msg.MinSelfDelegation,
		height,
	))
	if err != nil {
		return err
	}

	// Save validator description
	desc, err := bstakingcommon.GetValidatorDescription(
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

// handleMsgBeginRedelegate handles a MsgBeginRedelegate storing the data inside the database
func handleMsgBeginRedelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate,
	client stakingtypes.QueryClient, db *bigdipperdb.Db,
) error {
	redelegation, err := staking.StoreRedelegationFromMessage(tx, index, msg, db)
	if err != nil {
		return err
	}

	// When the time expires, update the delegations and delete this redelegation
	time.AfterFunc(time.Until(redelegation.CompletionTime),
		bstakingcommon.RefreshDelegations(tx.Height, msg.DelegatorAddress, client, db))
	time.AfterFunc(time.Until(redelegation.CompletionTime),
		bstakingcommon.DeleteRedelegation(*redelegation, db))

	// Update the current delegations
	return bstakingcommon.UpdateDelegationsAndReplaceExisting(tx.Height, msg.DelegatorAddress, client, db)
}

// handleMsgUndelegate handles a MsgUndelegate storing the data inside the database
func handleMsgUndelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate,
	stakingClient stakingtypes.QueryClient, bankClient banktypes.QueryClient, db *bigdipperdb.Db,
) error {
	delegation, err := staking.StoreUnbondingDelegationFromMessage(tx, index, msg, db)
	if err != nil {
		return err
	}

	// When timer expires update the delegations, update the user balance and remove the unbonding delegation
	time.AfterFunc(time.Until(delegation.CompletionTimestamp),
		bstakingcommon.RefreshDelegations(tx.Height, msg.DelegatorAddress, stakingClient, db))
	time.AfterFunc(time.Until(delegation.CompletionTimestamp),
		bbankcommon.RefreshBalance(msg.DelegatorAddress, bankClient, db))
	time.AfterFunc(time.Until(delegation.CompletionTimestamp),
		bstakingcommon.DeleteUnbondingDelegation(*delegation, db))

	// Update the current delegations
	return bstakingcommon.UpdateDelegationsAndReplaceExisting(tx.Height, msg.DelegatorAddress, stakingClient, db)
}
