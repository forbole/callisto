package staking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *stakingtypes.MsgCreateValidator:
		return m.handleMsgCreateValidator(tx.Height, cosmosMsg)

	case *stakingtypes.MsgEditValidator:
		return m.handleEditValidator(tx.Height, cosmosMsg)

	// update validators statuses, voting power
	// and proposals validators satatus snapshots
	// when there is a voting power change
	case *stakingtypes.MsgDelegate:
		return m.UpdateValidatorStatuses()

	case *stakingtypes.MsgBeginRedelegate:
		return m.UpdateValidatorStatuses()

	case *stakingtypes.MsgUndelegate:
		return m.UpdateValidatorStatuses()

	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgCreateValidator handles properly a MsgCreateValidator instance by
// saving into the database all the data associated to such validator
func (m *Module) handleMsgCreateValidator(height int64, msg *stakingtypes.MsgCreateValidator) error {
	err := m.RefreshValidatorInfos(height, msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing validator from MsgCreateValidator: %s", err)
	}
	return nil
}

// handleEditValidator handles MsgEditValidator utils, updating the validator info
func (m *Module) handleEditValidator(height int64, msg *stakingtypes.MsgEditValidator) error {
	err := m.RefreshValidatorInfos(height, msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf("error while refreshing validator from MsgEditValidator: %s", err)
	}

	return nil
}
