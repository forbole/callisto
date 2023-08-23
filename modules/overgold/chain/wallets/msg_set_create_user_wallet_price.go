package wallets

import (
	wallets "git.ooo.ua/vipcoin/chain/x/wallets/types"
	juno "github.com/forbole/juno/v3/types"
)

// handleMsgPayments allows to properly handle a MsgSetState
func (m *Module) handleMsgSetCreateUserWalletPrice(tx *juno.Tx, _ int, msg *wallets.MsgSetCreateUserWalletPrice) error {
	return m.walletsRepo.SaveSetCreateUserWalletPrice(msg, tx.TxHash)
}
