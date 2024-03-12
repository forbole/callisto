package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgTransferFromUser allows to properly handle a transfer from user message.
func (m *Module) handleMsgTransferFromUser(tx *juno.Tx, _ int, msg *types.MsgTransferFromUser) error {
	return m.stakeRepo.InsertMsgTransferFromUser(tx.TxHash, types.MsgTransferFromUser{
		Creator: msg.Creator,
		Amount:  msg.Amount,
		Address: msg.Address,
	})
}
