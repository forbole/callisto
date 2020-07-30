package database

import (
	"fmt"
	"time"

	api "github.com/forbole/bdjuno/x/pricefeed/apiTypes"
)

// SaveTokensPrice allows to save for the given height the given total amount of coins
func (db BigDipperDb) SaveTokensPrice(pricefeeds api.MarketTickers, timestamp time.Time) error {
	query := `INSERT INTO token_price(denom,price,market_cap,timestamp) VALUES`
	var param []interface{}
	for i, pricefeed := range pricefeeds {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, pricefeed.ID, pricefeed.CurrentPrice, pricefeed.MarketCap, timestamp)
	}
	query = query[:len(query)-1] // Remove trailing ","
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}
