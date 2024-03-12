package stake

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgDistributeRewards allows to properly handle a stake distribute rewards message
func (m *Module) handleMsgDistributeRewards(tx *juno.Tx, _ int, msg *types.MsgDistributeRewards) error {
	return m.stakeRepo.InsertMsgDistributeRewards(tx.TxHash, types.MsgDistributeRewards{
		Creator: msg.Creator,
	})
}
