package allowed

import (
	"fmt"

	"git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch allowedMsg := msg.(type) {
	case *types.MsgCreateAddresses:
		return m.handleMsgCreateAddresses(tx, index, allowedMsg)
	case *types.MsgDeleteByAddresses:
		return m.handleMsgDeleteByAddresses(tx, index, allowedMsg)
	case *types.MsgDeleteByID:
		return m.handleMsgDeleteByID(tx, index, allowedMsg)
	case *types.MsgUpdateAddresses:
		return m.handleMsgUpdateAddresses(tx, index, allowedMsg)
	default:
		return fmt.Errorf("unrecognized %s message type: %T", m.Name(), allowedMsg)
	}
}
