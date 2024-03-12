package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgSellCancel allows to properly handle a stake sell cancel message
func (m *Module) handleMsgSellCancel(tx *juno.Tx, _ int, msg *types.MsgMsgCancelSell) error {
	return m.stakeRepo.InsertMsgSellCancel(tx.TxHash, types.MsgMsgCancelSell{
		Creator: msg.Creator,
		Amount:  msg.Amount,
	})
}
