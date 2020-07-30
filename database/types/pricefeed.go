package types

import "time"

// TokenPriceRow represent a row of the table token_price in the database
type TokenPriceRow struct {
	Denom     string    `db:"denom"`
	Price     float64   `db:"price"`
	MarketCap int64     `db:"market_cap"`
	Timestamp time.Time `db:"timestamp"`
}

// NewTokenPriceRow allows to easily create a new NewTokenPriceRow
func NewTokenPriceRow(denom string, currentPrice float64, marketCap int64, timestamp time.Time) TokenPriceRow {
	return TokenPriceRow{
		Denom:     denom,
		Price:     currentPrice,
		MarketCap: marketCap,
		Timestamp: timestamp,
	}
}

// Equals return true if u and v represent the same row
func (u TokenPriceRow) Equals(v TokenPriceRow) bool {
	return u.Denom == v.Denom &&
		u.Price == v.Price &&
		u.MarketCap == v.MarketCap &&
		u.Timestamp.Equal(v.Timestamp)
}
