package staking

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return m.handleMsgCreateValidator(tx.Height, cosmosMsg)

	case *stakingtypes.MsgEditValidator:
		return m.handleEditValidator(tx.Height, cosmosMsg)

	case *stakingtypes.MsgDelegate:
		return m.storeDelegationFromMessage(tx.Height, cosmosMsg)

	case *stakingtypes.MsgBeginRedelegate:
		return m.handleMsgBeginRedelegate(tx, index, cosmosMsg)

	case *stakingtypes.MsgUndelegate:
		return m.handleMsgUndelegate(tx, index, cosmosMsg)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func (m *Module) handleMsgCreateValidator(height int64, msg *stakingtypes.MsgCreateValidator) error {
	err := m.refreshValidatorInfos(height, msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing validator from MsgCreateValidator: %s", err)
	}

	// Get the first self delegation
	delegations, err := m.getValidatorDelegations(height, msg.ValidatorAddress)
	if err != nil {
		return nil
	}

	return m.db.SaveDelegations(delegations)
}

// handleEditValidator handles MsgEditValidator utils, updating the validator info and commission
func (m *Module) handleEditValidator(height int64, msg *stakingtypes.MsgEditValidator) error {
	err := m.refreshValidatorInfos(height, msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing validator from MsgEditValidator: %s", err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgBeginRedelegate handles a MsgBeginRedelegate storing the data inside the database
func (m *Module) handleMsgBeginRedelegate(tx *juno.Tx, index int, msg *stakingtypes.MsgBeginRedelegate) error {
	_, err := m.storeRedelegationFromMessage(tx, index, msg)
	if err != nil {
		return fmt.Errorf("error while storing redelegation from message: %s", err)
	}

	// Update the current delegations
	return m.refreshDelegatorDelegations(tx.Height, msg.DelegatorAddress)
}

// handleMsgUndelegate handles a MsgUndelegate storing the data inside the database
func (m *Module) handleMsgUndelegate(tx *juno.Tx, index int, msg *stakingtypes.MsgUndelegate) error {
	_, err := m.storeUnbondingDelegationFromMessage(tx, index, msg)
	if err != nil {
		return fmt.Errorf("error while storing unbonding delegation from message: %s", err)
	}

	// Update the current delegations
	return m.refreshDelegatorDelegations(tx.Height, msg.DelegatorAddress)
}
