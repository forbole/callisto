package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *Db) SaveInflation(inflation sdk.Dec, height int64) error {
	stmt := `
INSERT INTO inflation (value, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET value = excluded.value, 
        height = excluded.height 
WHERE inflation.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, inflation.String(), height)
	return err
}
