package allowed

import (
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgCreateAddresses allows to properly handle a MsgCreateAddresses
func (m *Module) handleMsgCreateAddresses(tx *juno.Tx, _ int, msg *allowed.MsgCreateAddresses) error {
	if err := m.allowedRepo.InsertToCreateAddresses(tx.TxHash, msg); err != nil {
		return err
	}

	return m.allowedRepo.InsertToAddresses(allowed.Addresses{
		Address: msg.Address,
		Creator: msg.Creator,
	})
}
