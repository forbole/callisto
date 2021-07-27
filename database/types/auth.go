package types

// AccountRow represents a single row inside the account table
type AccountRow struct {
	Address string `db:"address"`
	Details string `db:"details"`
}

// NewAccountRow allows to easily build a new AccountRow
func NewAccountRow(address string,details string) AccountRow {
	return AccountRow{
		Address: address,
		Details: details,
	}
}

// Equal tells whether a and b contain the same data
func (a AccountRow) Equal(b AccountRow) bool {
	return a.Address == b.Address && a.Details==b.Details
}
