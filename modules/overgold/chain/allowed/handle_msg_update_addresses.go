package allowed

import (
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgUpdateAddresses allows to properly handle a MsgUpdateAddresses
func (m *Module) handleMsgUpdateAddresses(tx *juno.Tx, _ int, msg *allowed.MsgUpdateAddresses) error {
	if err := m.allowedRepo.InsertToUpdateAddresses(tx.TxHash, msg); err != nil {
		return err
	}

	return m.allowedRepo.UpdateAddresses(allowed.Addresses{
		Id:      msg.Id,
		Address: msg.Address,
		Creator: msg.Creator,
	})
}
