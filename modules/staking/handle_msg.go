package staking

import (
	"time"

	"github.com/forbole/bdjuno/database"
	bankutils "github.com/forbole/bdjuno/modules/bank/utils"
	stakingutils "github.com/forbole/bdjuno/modules/staking/utils"
	"github.com/forbole/bdjuno/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to handle the different utils related to the staking module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg, stakingClient stakingtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
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
		return stakingutils.StoreDelegationFromMessage(tx.Height, cosmosMsg, stakingClient, db)

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
	height int64, msg *stakingtypes.MsgCreateValidator, cdc codec.Marshaler, db *database.Db,
) error {
	err := stakingutils.StoreValidatorFromMsgCreateValidator(height, msg, cdc, db)
	if err != nil {
		return err
	}

	// Save validator description
	description, err := stakingutils.ConvertValidatorDescription(
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

// handleEditValidator handles MsgEditValidator utils, updating the validator info and commission
func handleEditValidator(
	height int64, msg *stakingtypes.MsgEditValidator, db *database.Db,
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
	desc, err := stakingutils.ConvertValidatorDescription(
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
	client stakingtypes.QueryClient, db *database.Db,
) error {
	redelegation, err := stakingutils.StoreRedelegationFromMessage(tx, index, msg, db)
	if err != nil {
		return err
	}

	// When the time expires, update the delegations and delete this redelegation
	time.AfterFunc(time.Until(redelegation.CompletionTime),
		stakingutils.RefreshDelegations(tx.Height, msg.DelegatorAddress, client, db))
	time.AfterFunc(time.Until(redelegation.CompletionTime),
		stakingutils.DeleteRedelegation(*redelegation, db))

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(tx.Height, msg.DelegatorAddress, client, db)
}

// handleMsgUndelegate handles a MsgUndelegate storing the data inside the database
func handleMsgUndelegate(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate,
	stakingClient stakingtypes.QueryClient, bankClient banktypes.QueryClient, db *database.Db,
) error {
	delegation, err := stakingutils.StoreUnbondingDelegationFromMessage(tx, index, msg, db)
	if err != nil {
		return err
	}

	// When timer expires update the delegations, update the user balance and remove the unbonding delegation
	time.AfterFunc(time.Until(delegation.CompletionTimestamp),
		stakingutils.RefreshDelegations(tx.Height, msg.DelegatorAddress, stakingClient, db))
	time.AfterFunc(time.Until(delegation.CompletionTimestamp),
		bankutils.RefreshBalance(msg.DelegatorAddress, bankClient, db))
	time.AfterFunc(time.Until(delegation.CompletionTimestamp),
		stakingutils.DeleteUnbondingDelegation(*delegation, db))

	// Update the current delegations
	return stakingutils.UpdateDelegationsAndReplaceExisting(tx.Height, msg.DelegatorAddress, stakingClient, db)
}
