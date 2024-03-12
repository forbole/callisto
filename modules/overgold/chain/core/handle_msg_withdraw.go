package core

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgWithdraw allows to properly handle a MsgWithdraw
func (m *Module) handleMsgWithdraw(tx *juno.Tx, _ int, msg *types.MsgWithdraw) error {
	return m.coreRepo.InsertMsgWithdraw(tx.TxHash, types.MsgWithdraw{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Denom:   msg.Denom,
		Address: msg.Address,
	})
}
