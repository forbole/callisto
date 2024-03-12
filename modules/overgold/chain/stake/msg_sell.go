package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgSell allows to properly handle a stake sell message
func (m *Module) handleMsgSell(tx *juno.Tx, _ int, msg *types.MsgSellRequest) error {
	return m.stakeRepo.InsertMsgSell(tx.TxHash, types.MsgSellRequest{
		Creator: msg.Creator,
		Amount:  msg.Amount,
	})
}
