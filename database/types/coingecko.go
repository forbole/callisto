package types

// TokenPriceRow represent a row of the table tokenprice in the database
type TokenPriceRow struct {
	Denom        string  `db:"denom"`
	Currentprice float64 `db:"current_price"`
	Marketcap    int64   `db:"market_cap"`
	Height       int64   `db:"height"`
}

// NewTokenPriceRow allows to easily create a new NewTokenPriceRow
func NewTokenPriceRow(denom string, currentprice float64, marketcap int64, height int64) TokenPriceRow {
	return TokenPriceRow{
		Denom:        denom,
		Currentprice: currentprice,
		Marketcap:    marketcap,
		Height:       height,
	}
}

// Equals return true if u and v represent the same row
func (u TokenPriceRow) Equals(v TokenPriceRow) bool {
	return u.Denom == v.Denom &&
		u.Currentprice == v.Currentprice &&
		u.Marketcap == v.Marketcap &&
		u.Height == v.Height
}
