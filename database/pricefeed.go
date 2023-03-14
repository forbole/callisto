package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"

	"github.com/forbole/bdjuno/v3/types"

	"github.com/lib/pq"
)

// GetTokensPriceID returns the slice of price ids for all tokens stored in db
func (db *Db) GetTokensPriceID() ([]string, error) {
	query := `SELECT * FROM token_unit`

	var dbUnits []dbtypes.TokenUnitRow
	err := db.Sqlx.Select(&dbUnits, query)
	if err != nil {
		return nil, err
	}

	var units []string
	for _, unit := range dbUnits {
		if unit.PriceID.String != "" {
			units = append(units, unit.PriceID.String)
		}
	}

	return units, nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveToken allows to save the given token details
func (db *Db) SaveToken(token types.Token) error {
	query := `INSERT INTO token (name) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(query, token.Name)
	if err != nil {
		return err
	}

	query = `INSERT INTO token_unit (token_name, denom, exponent, aliases, price_id) VALUES `
	var params []interface{}

	for i, unit := range token.Units {
		ui := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", ui+1, ui+2, ui+3, ui+4, ui+5)
		params = append(params, token.Name, unit.Denom, unit.Exponent, pq.StringArray(unit.Aliases),
			dbtypes.ToNullString(unit.PriceID))
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err = db.SQL.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while saving token: %s", err)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// SaveTokensPrices allows to save the given prices as the most updated ones
func (db *Db) SaveTokensPrices(prices []types.TokenPrice) error {
	if len(prices) == 0 {
		return nil
	}

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

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving tokens prices: %s", err)
	}

	return nil
}

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

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while storing tokens price history: %s", err)
	}

	return nil
}
