package types

// NewCommunityPoolRow represents a single row inside the total_supply table
type CommunityPoolRow struct {
	Coins  *DbDecCoins `db:"coins"`
	Height int64       `db:"height"`
}

// NewCommunityPoolRow allows to easily create a new NewCommunityPoolRow
func NewCommunityPoolRow(coins DbDecCoins, height int64) CommunityPoolRow {
	return CommunityPoolRow{
		Coins:  &coins,
		Height: height,
	}
}

// Equals return true if one CommunityPoolRow representing the same row as the original one
func (v CommunityPoolRow) Equals(w CommunityPoolRow) bool {
	return v.Coins.Equal(w.Coins) &&
		v.Height == w.Height
}
