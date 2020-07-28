package types

// TokenPriceRow represent a row of the table token_price in the database
type TokenPriceRow struct {
	Denom         string
	Current_price float64
	Market_cap    int64
	Height        int64
}

// NewTokenPriceRow allows to easily create a new NewTokenPriceRow
func NewTokenPriceRow(denom string, current_price float64, market_cap int64, height int64) TokenPriceRow {
	return TokenPriceRow{
		Denom:         denom,
		Current_price: current_price,
		Market_cap:    market_cap,
		Height:        height,
	}
}

// Equals return true if one NewTokenPriceRow representing the same row as the original one
func (u TokenPriceRow) Equals(v TokenPriceRow) bool {
	return u.Denom == v.Denom &&
		u.Current_price == v.Current_price &&
		u.Market_cap == v.Market_cap &&
		u.Height == v.Height
}
