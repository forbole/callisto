package forbolex

import (
	"fmt"

	"github.com/lib/pq"

	"github.com/forbole/bdjuno/modules/common/pricefeed"
	"github.com/forbole/bdjuno/types"
)

var (
	_ pricefeed.DB = &Db{}
)

// SaveToken implements pricefeed.DB
func (db *Db) SaveToken(token types.Token) error {
	query := `INSERT INTO token (name) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query, token.Name)
	if err != nil {
		return err
	}

	query = `INSERT INTO token_unit (token_name, denom, exponent, aliases) VALUES `
	var params []interface{}

	for i, unit := range token.Units {
		ui := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", ui+1, ui+2, ui+3, ui+4)
		params = append(params, token.Name, unit.Denom, unit.Exponent, pq.StringArray(unit.Aliases))
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(query, params...)
	return err
}

// GetTokenUnits implements pricefeed.DB
func (db *Db) GetTokenUnits() ([]string, error) {
	query := `SELECT denom FROM token_unit`

	var names pq.StringArray
	err := db.Sqlx.Select(&names, query)
	if err != nil {
		return nil, err
	}

	return names, nil
}

// SaveTokensPrices implements pricefeed.DB
func (db *Db) SaveTokensPrices(prices []types.TokenPrice) error {
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
	return err
}
