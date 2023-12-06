package core

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch coreMsg := msg.(type) {
	case *types.MsgIssue:
		return m.handleMsgIssue(tx, index, coreMsg)
	case *types.MsgWithdraw:
		return m.handleMsgWithdraw(tx, index, coreMsg)
	case *types.MsgSend:
		return m.handleMsgSend(tx, index, coreMsg)
	default:
		return nil
	}
}
