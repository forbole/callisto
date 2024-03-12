package feeexcluder

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgDeleteAddress allows to properly handle a message
func (m *Module) handleMsgDeleteAddress(tx *juno.Tx, _ int, msg *types.MsgDeleteAddress) error {
	if err := m.feeexcluderRepo.InsertToMsgDeleteAddress(tx.TxHash, *msg); err != nil {
		return err
	}

	return m.feeexcluderRepo.DeleteAddress(nil, msg.Id)
}
