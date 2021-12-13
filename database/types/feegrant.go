package types

// FeeAllowanceRow represents a single row inside the fee_grant_allowance table
type FeeAllowanceRow struct {
	ID        uint64 `db:"id"`
	Grantee   string `db:"grantee_address"`
	Granter   string `db:"granter_address"`
	Allowance string `db:"allowance"`
	Height    int64  `db:"height"`
}
