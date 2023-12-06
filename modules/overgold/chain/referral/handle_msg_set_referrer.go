package referral

import (
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	referral "git.ooo.ua/vipcoin/ovg-chain/x/referral/types"
	juno "github.com/forbole/juno/v5/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// handleMsgSetReferrer allows to properly handle a MsgSetReferrer
func (m *Module) handleMsgSetReferrer(tx *juno.Tx, index int, msg *referral.MsgSetReferrer) error {
	msgs, err := m.referralRepo.GetAllMsgSetReferrer(filter.NewFilter().SetCondition(filter.ConditionAND).
		SetArgument(db.FieldTxHash, tx.TxHash).
		SetArgument(db.FieldCreator, msg.Creator))
	if err != nil {
		if !errors.As(err, &errs.NotFound{}) {
			return err
		}
	}
	if len(msgs) > 0 {
		return errs.AlreadyExists{What: "msg set referrer, referrer address: " + msg.ReferrerAddress}
	}

	return m.referralRepo.InsertMsgSetReferrer(tx.TxHash, referral.MsgSetReferrer{
		Creator:         msg.Creator,
		ReferrerAddress: msg.ReferrerAddress,
		ReferralAddress: msg.ReferralAddress,
	})
}
