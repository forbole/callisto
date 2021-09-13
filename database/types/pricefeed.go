package types

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type TokenUnitRow struct {
	TokenName string         `db:"token_name"`
	Denom     string         `db:"denom"`
	Exponent  int            `db:"exponent"`
	Aliases   pq.StringArray `db:"aliases"`
	PriceID   sql.NullString `db:"price_id"`
}

type TokenRow struct {
	Name       string `db:"name"`
	TradedUnit string `db:"traded_unit"`
}

// --------------------------------------------------------------------------------------------------------------------

// TokenPriceRow represent a row of the table token_price in the database
type TokenPriceRow struct {
	ID        string    `db:"id"`
	Name      string    `db:"unit_name"`
	Price     float64   `db:"price"`
	MarketCap int64     `db:"market_cap"`
	Timestamp time.Time `db:"timestamp"`
}

// NewTokenPriceRow allows to easily create a new NewTokenPriceRow
func NewTokenPriceRow(name string, currentPrice float64, marketCap int64, timestamp time.Time) TokenPriceRow {
	return TokenPriceRow{
		Name:      name,
		Price:     currentPrice,
		MarketCap: marketCap,
		Timestamp: timestamp,
	}
}

// Equals return true if u and v represent the same row
func (u TokenPriceRow) Equals(v TokenPriceRow) bool {
	return u.Name == v.Name &&
		u.Price == v.Price &&
		u.MarketCap == v.MarketCap &&
		u.Timestamp.Equal(v.Timestamp)
}
