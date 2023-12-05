package referral

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/referral/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// toMsgSetReferrerDomain - mapping func to a domain model.
func toMsgSetReferrerDomain(m db.DbReferralSetReferrer) types.MsgSetReferrer {
	return types.MsgSetReferrer{
		Creator:         m.Creator,
		ReferrerAddress: m.ReferralAddress,
		ReferralAddress: m.ReferralAddress,
	}
}

// toMsgSetReferrerDomainList - mapping func to a domain list.
func toMsgSetReferrerDomainList(m []db.DbReferralSetReferrer) []types.MsgSetReferrer {
	res := make([]types.MsgSetReferrer, 0, len(m))
	for _, msg := range m {
		res = append(res, toMsgSetReferrerDomain(msg))
	}

	return res
}

// toMsgSetReferrerDatabase - mapping func to a database model.
func toMsgSetReferrerDatabase(hash string, m types.MsgSetReferrer) db.DbReferralSetReferrer {
	return db.DbReferralSetReferrer{
		TxHash:          hash,
		Creator:         m.Creator,
		ReferrerAddress: m.ReferrerAddress,
		ReferralAddress: m.ReferralAddress,
	}
}
