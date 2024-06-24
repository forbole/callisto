package staking

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	juno "github.com/forbole/juno/v6/types"

	"github.com/forbole/callisto/v4/utils"
)

var msgFilter = map[string]bool{
	"/cosmos.staking.v1beta1.MsgCreateValidator": true,
	"/cosmos.staking.v1beta1.MsgEditValidator":   true,
	"/cosmos.staking.v1beta1.MsgDelegate":        true,
	"/cosmos.staking.v1beta1.MsgUndelegate":      true,
	"/cosmos.staking.v1beta1.MsgBeginRedelegate": true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(_ int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	log.Debug().Str("module", "staking").Str("hash", tx.TxHash).Uint64("height", tx.Height).Msg(fmt.Sprintf("handling staking message %s", msg.GetType()))

	switch msg.GetType() {
	case "/cosmos.staking.v1beta1.MsgCreateValidator":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &stakingtypes.MsgCreateValidator{})
		return m.handleMsgCreateValidator(int64(tx.Height), cosmosMsg)

	case "/cosmos.staking.v1beta1.MsgEditValidator":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &stakingtypes.MsgEditValidator{})
		return m.handleEditValidator(int64(tx.Height), cosmosMsg)

	// update validators statuses, voting power
	// and proposals validators satatus snapshots
	// when there is a voting power change
	case "/cosmos.staking.v1beta1.MsgDelegate":
		return m.UpdateValidatorStatuses()

	case "/cosmos.staking.v1beta1.MsgBeginRedelegate":
		return m.UpdateValidatorStatuses()

	case "/cosmos.staking.v1beta1.MsgUndelegate":
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
