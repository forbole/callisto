package database

import (
	"fmt"

	"github.com/lib/pq"

	"github.com/forbole/bdjuno/x/pricefeed/types"
)

// GetTradedNames returns the slice of all the names of the different tokens that can traded on exchanges
func (db *BigDipperDb) GetTradedNames() ([]string, error) {
	query := `SELECT traded_unit FROM token`

	var names pq.StringArray
	err := db.Sqlx.Select(&names, query)
	if err != nil {
		return nil, err
	}

	return names, nil
}

// SaveTokensPrices allows to save the given tickers associating them to the given timestamp
func (db *BigDipperDb) SaveTokensPrices(tickers types.MarketTickers) error {
	query := `INSERT INTO token_price (name, price, market_cap, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range tickers {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, ticker.Symbol, ticker.CurrentPrice, ticker.MarketCap, ticker.LastUpdated)
	}

	query = query[:len(query)-1] // Remove trailing ","
	_, err := db.Sql.Exec(query, param...)
	return err
}
