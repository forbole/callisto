package bank

import (
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgSend allows to properly handle a MsgSend
func (m *Module) handleMsgSend(tx *juno.Tx, _ int, msg *bank.MsgSend) error {
	return m.bankRepo.InsertMsgSend(tx.TxHash, bank.MsgSend{
		FromAddress: msg.FromAddress,
		ToAddress:   msg.ToAddress,
		Amount:      msg.Amount,
	})
}
