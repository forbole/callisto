package referral

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch referralMsg := msg.(type) {
	case *types.MsgSetReferrer:
		return m.handleMsgSetReferrer(tx, index, referralMsg)
	default:
		return nil
	}
}
