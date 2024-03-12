package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgBuy allows to properly handle a stake buy message
func (m *Module) handleMsgBuy(tx *juno.Tx, _ int, msg *types.MsgBuyRequest) error {
	return m.stakeRepo.InsertMsgBuy(tx.TxHash, types.MsgBuyRequest{
		Creator: msg.Creator,
		Amount:  msg.Amount,
	})
}
