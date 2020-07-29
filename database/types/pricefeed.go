package types

// TokenPriceRow represent a row of the table token_price in the database
type TokenPriceRow struct {
	Denom     string  `db:"denom"`
	Price     float64 `db:"current_price"`
	MarketCap int64   `db:"market_cap"`
	Height    int64   `db:"height"`
}

// NewTokenPriceRow allows to easily create a new NewTokenPriceRow
func NewTokenPriceRow(denom string, currentPrice float64, marketCap int64, height int64) TokenPriceRow {
	return TokenPriceRow{
		Denom:     denom,
		Price:     currentPrice,
		MarketCap: marketCap,
		Height:    height,
	}
}

// Equals return true if u and v represent the same row
func (u TokenPriceRow) Equals(v TokenPriceRow) bool {
	return u.Denom == v.Denom &&
		u.Price == v.Price &&
		u.MarketCap == v.MarketCap &&
		u.Height == v.Height
}
