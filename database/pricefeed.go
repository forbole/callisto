package database

import (
	"fmt"
	"time"

	pricefeedtypes "github.com/forbole/bdjuno/x/pricefeed/types"
)

// SaveTokensPrices allows to save the given tickers associating them to the given timestamp
func (db BigDipperDb) SaveTokensPrices(tickers pricefeedtypes.MarketTickers, timestamp time.Time) error {
	query := `INSERT INTO token_price (denom, price, market_cap, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range tickers {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, ticker.ID, ticker.CurrentPrice, ticker.MarketCap, timestamp)
	}

	query = query[:len(query)-1] // Remove trailing ","
	_, err := db.Sql.Exec(query, param...)
	return err
}
