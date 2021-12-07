package types

// FeeAllowanceRow represents a single row inside the fee_grant_allowance table
type FeeAllowanceRow struct {
	Grantee   string `db:"grantee_address"`
	Granter   string `db:"granter_address"`
	Allowance string `db:"allowance"`
	Height    int64  `db:"height"`
}

// NewFeeAllowanceRow allows to easily create a new FeeAllowanceRow
func NewFeeAllowanceRow(
	grantee string,
	granter string,
	allowance string,
	height int64,
) FeeAllowanceRow {
	return FeeAllowanceRow{
		Grantee:   grantee,
		Granter:   granter,
		Allowance: allowance,
		Height:    height,
	}
}

// Equals return true if two FeeAllowanceRow are the same
func (w FeeAllowanceRow) Equals(v FeeAllowanceRow) bool {
	return w.Grantee == v.Grantee &&
		w.Granter == v.Granter &&
		w.Allowance == w.Allowance &&
		w.Height == v.Height
}
