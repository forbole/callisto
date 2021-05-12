package types

// SupplyRow represents a single row inside the "supply" table
type SupplyRow struct {
	OneRowID bool     `db:"one_row_id"`
	Coins    *DbCoins `db:"coins"`
	Height   int64    `db:"height"`
}

// NewSupplyRow allows to easily create a new NewSupplyRow
func NewSupplyRow(coins DbCoins, height int64) SupplyRow {
	return SupplyRow{
		OneRowID: true,
		Coins:    &coins,
		Height:   height,
	}
}

// Equals return true if one totalSupplyRow representing the same row as the original one
func (v SupplyRow) Equals(w SupplyRow) bool {
	return v.Coins.Equal(w.Coins) &&
		v.Height == w.Height
}
