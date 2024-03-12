package referral

import (
	referral "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgSetReferrer allows to properly handle a MsgSetReferrer
func (m *Module) handleMsgSetReferrer(tx *juno.Tx, _ int, msg *referral.MsgSetReferrer) error {
	return m.referralRepo.InsertMsgSetReferrer(tx.TxHash, referral.MsgSetReferrer{
		Creator:         msg.Creator,
		ReferrerAddress: msg.ReferrerAddress,
		ReferralAddress: msg.ReferralAddress,
	})
}
