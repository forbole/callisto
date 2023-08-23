package banking

import (
	"git.ooo.ua/vipcoin/chain/x/banking/types"
	juno "github.com/forbole/juno/v3/types"
)

// handleMsgSetRewardManagerAddress allows to properly handle a handleMsgSetRewardManagerAddress
func (m *Module) handleMsgSetRewardManagerAddress(tx *juno.Tx, index int, msg *types.MsgSetRewardManagerAddress) error {
	if err := m.bankingRepo.SaveMsgSetRewardMgrAddress(msg, tx.TxHash); err != nil {
		return err
	}

	return nil
}
