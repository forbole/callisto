package banking

import (
	"fmt"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	juno "github.com/forbole/juno/v2/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch bankingMsg := msg.(type) {
	case *types.MsgSystemTransfer:
		return m.handleMsgSystemTransfer(tx, index, bankingMsg)

	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, bankingMsg)
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
