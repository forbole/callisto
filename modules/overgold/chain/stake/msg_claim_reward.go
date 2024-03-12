package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgClaimReward allows to properly handle a stake claim reward message
func (m *Module) handleMsgClaimReward(tx *juno.Tx, _ int, msg *types.MsgClaimReward) error {
	return m.stakeRepo.InsertMsgClaimReward(tx.TxHash, types.MsgClaimReward{
		Creator: msg.Creator,
	})
}
