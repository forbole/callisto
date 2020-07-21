package types

// NewTotalSupplyRow represents a single row inside the total_supply table
type TotalSupplyRow struct {
	Coins  *DbCoins `db:"coins"`
	Height int64    `db:"height"`
}

// NewTotalSupplyRow allows to easily create a new NewTotalSupplyRow
func NewTotalSupplyRow(coins DbCoins, height int64) TotalSupplyRow {
	return TotalSupplyRow{
		Coins:  &coins,
		Height: height,
	}
}

// Equals return true if one totalSupplyRow representing the same row as the original one
func (v TotalSupplyRow) Equals(w TotalSupplyRow) bool {
	return v.Coins.Equal(w.Coins) &&
		v.Height == w.Height
}
