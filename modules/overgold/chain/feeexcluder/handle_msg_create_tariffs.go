package feeexcluder

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgCreateTariffs allows to properly handle a message
func (m *Module) handleMsgCreateTariffs(tx *juno.Tx, _ int, msg *types.MsgCreateTariffs) error {
	return m.feeexcluderRepo.InsertToMsgCreateTariffs(tx.TxHash, *msg)
}
