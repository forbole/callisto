package banking

import (
	"time"

	"git.ooo.ua/vipcoin/chain/x/banking/types"
	juno "github.com/forbole/juno/v2/types"
)

// handleMsgSystemRewardTransfer allows to properly handle a MsgSystemRewardTransfer
func (m *Module) handleMsgSystemRewardTransfer(tx *juno.Tx, index int, msg *types.MsgSystemRewardTransfer) error {
	if err := m.bankingRepo.SaveMsgSystemRewardTransfers(msg); err != nil {
		return err
	}

	time, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	transfer := &types.SystemTransfer{
		WalletFrom: msg.WalletFrom,
		WalletTo:   msg.WalletTo,
		BaseTransfer: types.BaseTransfer{
			Asset:     msg.Asset,
			Amount:    msg.Amount,
			Kind:      types.TRANSFER_KIND_SYSTEM_REWARD,
			Extras:    msg.Extras,
			Timestamp: time.Unix(),
			TxHash:    tx.TxHash,
		},
	}

	return m.bankingRepo.SaveSystemTransfers(transfer)
}
