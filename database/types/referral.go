package types

type (
	// DbReferralSetReferrer - table for storing referral set
	DbReferralSetReferrer struct {
		ID              uint64 `db:"id"`
		TxHash          string `db:"tx_hash"`
		Creator         string `db:"creator"`
		ReferrerAddress string `db:"referrer_address"`
		ReferralAddress string `db:"referral_address"`
	}
)
