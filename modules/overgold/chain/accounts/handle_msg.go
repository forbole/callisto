package accounts

import (
	"git.ooo.ua/vipcoin/chain/x/accounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
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
	case *types.MsgCreateAccount:
		return m.handleMsgCreateAccount(tx, index, accountMsg)
	case *types.MsgAddAffiliate:
		return m.handleMsgAddAffiliate(tx, index, accountMsg)
	default:
		return nil
	}
}
