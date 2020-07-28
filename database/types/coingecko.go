package types

// TokenPriceRow represent a row of the table tokenprice in the database
type TokenPriceRow struct {
	Denom        string
	Currentprice float64
	Marketcap    int64
	Height       int64
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

// Equals return true if one NewTokenPriceRow representing the same row as the original one
func (u TokenPriceRow) Equals(v TokenPriceRow) bool {
	return u.Denom == v.Denom &&
		u.Currentprice == v.Currentprice &&
		u.Marketcap == v.Marketcap &&
		u.Height == v.Height
}
