package types

// DistributionParamsRow represents a single row inside the distribution_params table
type DistributionParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// -------------------------------------------------------------------------------------------------------------------

// CommunityPoolRow represents a single row inside the total_supply table
type CommunityPoolRow struct {
	OneRowID bool        `db:"one_row_id"`
	Coins    *DbDecCoins `db:"coins"`
	Height   int64       `db:"height"`
}

// NewCommunityPoolRow allows to easily create a new CommunityPoolRow
func NewCommunityPoolRow(coins DbDecCoins, height int64) CommunityPoolRow {
	return CommunityPoolRow{
		OneRowID: true,
		Coins:    &coins,
		Height:   height,
	}
}

// Equals return true if one CommunityPoolRow representing the same row as the original one
func (v CommunityPoolRow) Equals(w CommunityPoolRow) bool {
	return v.Coins.Equal(w.Coins) &&
		v.Height == w.Height
}
