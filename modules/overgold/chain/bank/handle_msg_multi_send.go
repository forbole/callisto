package bank

import (
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgMultiSend allows to properly handle a MsgMultiSend
func (m *Module) handleMsgMultiSend(tx *juno.Tx, _ int, msg *bank.MsgMultiSend) error {
	return m.bankRepo.InsertMsgMultiSend(tx.TxHash, bank.MsgMultiSend{
		Inputs:  msg.Inputs,
		Outputs: msg.Outputs,
	})
}
