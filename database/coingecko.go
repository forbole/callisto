package database

import (
	"fmt"

	api "github.com/forbole/bdjuno/x/supply/coinGeckoTypes"
)

// SaveTokensPrice allows to save for the given height the given total amount of coins
func (db BigDipperDb) SaveTokensPrice(markets api.Markets, height int64) error {
	query := `INSERT INTO token_values(denom,current_price,market_cap,height) VALUES`
	var param []interface{}
	for i, market := range markets {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, market.Id, market.CurrentPrice, market.MarketCap,height)
	}
	query = query[:len(query)-1] // Remove trailing ","
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}
