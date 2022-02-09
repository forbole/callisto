package database

import (
	"fmt"

	"github.com/forbole/bdjuno/v2/types"
)

// SaveTokenPricesHistory stores the given prices as historic ones
func (db *Db) SaveTokenPricesHistory(prices []types.TokenPrice) error {
	if len(prices) == 0 {
		return nil
	}

	query := `INSERT INTO token_price_history (unit_name, price, market_cap, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range prices {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, ticker.UnitName, ticker.Price, ticker.MarketCap, ticker.Timestamp)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT ON CONSTRAINT unique_price_for_timestamp DO UPDATE 
	SET price = excluded.price,
	    market_cap = excluded.market_cap`

	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while storing tokens price history: %s", err)
	}

	return nil
}
