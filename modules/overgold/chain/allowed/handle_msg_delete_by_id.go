package allowed

import (
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgDeleteByID allows to properly handle a MsgDeleteByID
func (m *Module) handleMsgDeleteByID(tx *juno.Tx, _ int, msg *allowed.MsgDeleteByID) error {
	if err := m.allowedRepo.InsertToDeleteByID(tx.TxHash, msg); err != nil {
		return err
	}

	return m.allowedRepo.DeleteAddressesByID(msg.Id)
}
