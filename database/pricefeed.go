package database

import (
	"fmt"

	"github.com/forbole/bdjuno/types"

	"github.com/lib/pq"
)

// GetTokenUnits returns the slice of all the names of the different tokens units
func (db *Db) GetTokenUnits() ([]string, error) {
	query := `SELECT denom FROM token_unit`

	var names pq.StringArray
	err := db.Sqlx.Select(&names, query)
	if err != nil {
		return nil, err
	}

	return names, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveToken allows to save the given token details
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

// --------------------------------------------------------------------------------------------------------------------

// SaveTokensPrices allows to save the given prices as the most updated ones
func (db *Db) SaveTokensPrices(prices []types.TokenPrice) error {
	if len(prices) == 0 {
		return nil
	}

	err := db.saveUpToDateTokenPrices(prices)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date token prices: %s", err)
	}

	if db.IsStoreHistoricDataEnabled() {
		err = db.saveTokenPricesHistory(prices)
		if err != nil {
			return fmt.Errorf("error while storing historic token prices: %s", err)
		}
	}

	return nil
}

// saveUpToDateTokenPrices stores the given prices as the most up-to-date ones
func (db *Db) saveUpToDateTokenPrices(prices []types.TokenPrice) error {
	query := `INSERT INTO token_price (unit_name, price, market_cap, timestamp) VALUES`
	var param []interface{}

	for i, ticker := range prices {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, ticker.UnitName, ticker.Price, ticker.MarketCap, ticker.Timestamp)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (unit_name) DO UPDATE 
	SET price = excluded.price,
	    market_cap = excluded.market_cap,
	    timestamp = excluded.timestamp
WHERE token_price.timestamp <= excluded.timestamp`

	_, err := db.Sql.Exec(query, param...)
	return err
}

// saveTokenPricesHistory stores the given prices as historic ones
func (db *Db) saveTokenPricesHistory(prices []types.TokenPrice) error {
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
