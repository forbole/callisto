package staking

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	juno "github.com/desmos-labs/juno/types"
)

// storeValidatorFromMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func (m *Module) storeValidatorFromMsgCreateValidator(height int64, msg *stakingtypes.MsgCreateValidator) error {
	var pubKey cryptotypes.PubKey
	err := m.cdc.UnpackAny(msg.Pubkey, &pubKey)
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

	validator, err := m.convertValidator(height, stakingValidator)
	if err != nil {
		return fmt.Errorf("error while converting validator: %s", err)
	}

	desc, err := m.convertValidatorDescription(height, msg.ValidatorAddress, msg.Description)
	if err != nil {
		return fmt.Errorf("error while converting validator description: %s", err)
	}

	// Save the validator
	err = m.db.SaveValidatorsData([]types.Validator{validator})
	if err != nil {
		return err
	}

	// Save the description
	err = m.db.SaveValidatorDescription(desc)
	if err != nil {
		return err
	}

	// Save the first self-delegation
	err = m.db.SaveDelegations([]types.Delegation{
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
	return m.db.SaveValidatorCommission(types.NewValidatorCommission(
		msg.ValidatorAddress,
		&msg.Commission.Rate,
		&msg.MinSelfDelegation,
		height,
	))
}

// storeDelegationFromMessage handles a MsgDelegate and saves the delegation inside the database
func (m *Module) storeDelegationFromMessage(height int64, msg *stakingtypes.MsgDelegate) error {
	delegation, err := m.source.GetDelegation(height, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return err
	}

	return m.db.SaveDelegations([]types.Delegation{
		convertDelegationResponse(height, delegation),
	})
}

// storeRedelegationFromMessage handles a MsgBeginRedelegate by saving the redelegation inside the database,
// and returns the new redelegation instance
func (m *Module) storeRedelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate,
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

	return &redelegation, m.db.SaveRedelegations([]types.Redelegation{redelegation})
}

// storeUnbondingDelegationFromMessage handles a MsgUndelegate storing the new unbonding delegation inside the database,
// and returns the new unbonding delegation instance
func (m *Module) storeUnbondingDelegationFromMessage(
	tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate,
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

	return &delegation, m.db.SaveUnbondingDelegations([]types.UnbondingDelegation{delegation})
}
