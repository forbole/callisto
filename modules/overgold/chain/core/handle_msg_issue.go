package core

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/core/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgIssue allows to properly handle a MsgIssue
func (m *Module) handleMsgIssue(tx *juno.Tx, _ int, msg *types.MsgIssue) error {

	return m.coreRepo.InsertMsgIssue(tx.TxHash, types.MsgIssue{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Denom:   msg.Denom,
		Address: msg.Address,
	})
}
