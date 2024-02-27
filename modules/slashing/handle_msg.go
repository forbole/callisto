package slashing

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if cosmosMsg, ok := msg.(*slashingtypes.MsgUnjail); ok {
		return m.handleMsgUnjail(tx, cosmosMsg)
	}

	return nil
}

// handleMsgUnjail handles a MsgUnjail message by refreshing the validator info, status and voting power
func (m *Module) handleMsgUnjail(tx *juno.Tx, msg *slashingtypes.MsgUnjail) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return m.stakingModule.RefreshValidatorInfos(tx.Height, timestamp, msg.ValidatorAddr)
}
