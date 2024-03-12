package core

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgSend allows to properly handle a MsgSend
func (m *Module) handleMsgSend(tx *juno.Tx, _ int, msg *types.MsgSend) error {
	return m.coreRepo.InsertMsgSend(tx.TxHash, types.MsgSend{
		Creator: msg.Creator,
		From:    msg.From,
		To:      msg.To,
		Amount:  msg.Amount,
		Denom:   msg.Denom,
	})
}
