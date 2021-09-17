package database

import (
	"github.com/forbole/bdjuno/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *Db) SaveEmoneyInflation(inflation types.EmoneyInflation) error {
	stmt := `
INSERT INTO emoney_inflation (issuer, denom, rate, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET issuer = excluded.issuer, 
		denom = excluded.denom,
		rate = excluded.rate,
		height = height.rate
WHERE emoney_inflation.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, inflation.Issuer, inflation.Rate, inflation.Rate.String(), inflation.Height)
	if err != nil {
		return err
	}
	return nil
}
