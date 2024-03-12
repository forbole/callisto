package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgTransferToUser allows to properly handle a transfer to user message.
func (m *Module) handleMsgTransferToUser(tx *juno.Tx, _ int, msg *types.MsgTransferToUser) error {
	return m.stakeRepo.InsertMsgTransferToUser(tx.TxHash, types.MsgTransferToUser{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Address: msg.Address,
	})
}
