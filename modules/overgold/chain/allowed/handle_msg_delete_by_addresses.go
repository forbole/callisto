package allowed

import (
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgDeleteByAddresses allows to properly handle a MsgDeleteByAddresses
func (m *Module) handleMsgDeleteByAddresses(tx *juno.Tx, _ int, msg *allowed.MsgDeleteByAddresses) error {
	if err := m.allowedRepo.InsertToDeleteByAddresses(tx.TxHash, msg); err != nil {
		return err
	}

	return m.allowedRepo.DeleteAddressesByAddress(msg.Address...)
}
