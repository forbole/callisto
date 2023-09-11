package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	juno "github.com/forbole/juno/v4/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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

	if _, ok := msg.(*distrtypes.MsgFundCommunityPool); ok {
		return m.updateCommunityPool(tx.Height)
	}
	return nil
}
