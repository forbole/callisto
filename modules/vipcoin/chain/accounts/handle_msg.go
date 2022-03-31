package accounts

import (
	"fmt"

	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	juno "github.com/forbole/juno/v2/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch accountMsg := msg.(type) {
	case *types.MsgSetKinds:
		return m.handleMsgSetKinds(tx, index, accountMsg)
	case *types.MsgSetState:
		return m.handleMsgSetState(tx, index, accountMsg)
	case *types.MsgSetAffiliateAddress:
		return m.handleMsgSetAffiliateAddress(tx, index, accountMsg)
	case *types.MsgAccountMigrate:
		return m.handleMsgAccountMigrate(tx, index, accountMsg)
	case *types.MsgSetAffiliateExtra:
		return m.handleMsgSetAffiliateExtra(tx, index, accountMsg)
	case *types.MsgRegisterUser:
		return m.handleMsgRegisterUser(tx, index, accountMsg)
	case *types.MsgSetExtra:
		return m.handleMsgSetExtra(tx, index, accountMsg)
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, accountMsg)
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
