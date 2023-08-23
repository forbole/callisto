package wallets

import (
	typeswallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch walletMsg := msg.(type) {
	case *typeswallets.MsgSetWalletState:
		return m.handleMsgSetStates(tx, index, walletMsg)
	case *typeswallets.MsgCreateWallet:
		return m.handleMsgCreateWallet(tx, index, walletMsg)
	case *typeswallets.MsgSetDefaultWallet:
		return m.handleMsgSetDefaultWallet(tx, index, walletMsg)
	case *typeswallets.MsgSetExtra:
		return m.handleMsgSetExtra(tx, index, walletMsg)
	case *typeswallets.MsgCreateWalletWithBalance:
		return m.MsgCreateWalletWithBalance(tx, index, walletMsg)
	case *typeswallets.MsgSetWalletKind:
		return m.handleMsgSetKind(tx, index, walletMsg)
	case *typeswallets.MsgSetCreateUserWalletPrice:
		return m.handleMsgSetCreateUserWalletPrice(tx, index, walletMsg)
	default:
		return nil
	}
}
