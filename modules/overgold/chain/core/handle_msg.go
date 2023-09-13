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

	// switch coreMsg := msg.(type) {
	switch msg.(type) {
	case *types.MsgIssue:
		return nil
	case *types.MsgWithdraw:
		return nil
	default:
		return nil
	}
}
